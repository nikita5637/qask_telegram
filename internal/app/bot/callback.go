package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"qask_telegram/internal/app/model"
	"qask_telegram/internal/app/router"
	"qask_telegram/internal/app/store"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type callBackQueryHandler struct {
	bot    *tgbotapi.BotAPI
	config *Config
	logger *logrus.Logger
	router *router.Router
	store  store.Store
}

func newCallBackQueryHandler(bot *tgbotapi.BotAPI, config *Config, logger *logrus.Logger, store store.Store) *callBackQueryHandler {
	cH := &callBackQueryHandler{
		bot:    bot,
		config: config,
		logger: logger,
		router: router.NewRouter(logger),
		store:  store,
	}

	cH.configureRouter()

	return cH
}

func (h *callBackQueryHandler) configureRouter() {
	h.logger.Debugf("Configuring callback commands router ...")
	// Registering new routes (path, isPublic, handler)
	h.router.NewRoute("/register", false, h.handleRegisterUser())
	h.router.NewRoute("/profile", false, h.handleProfile())
	h.router.NewRoute("/getQuestion", false, h.handleGetQuestion())
	h.router.NewRoute("/showAnswer", false, h.handleShowAnswer())
	h.router.NewRoute("/showQuestion", false, h.handleShowQuestion())
	h.router.NewRoute("/showComment", false, h.handleShowComment())
	h.router.NewRoute("/getMathProblem", false, h.handleGetMathProblem())
	h.router.NewRoute("/settings", false, h.handleSettings())
	h.router.NewRoute("/subscribtions", false, h.handleSubscriptions())
	h.router.NewRoute("/back", false, h.handleBack())
	h.router.NewRoute("/sendReport", false, h.handleSendReport())
	h.router.NewRoute("/setFirstName", false, h.handleSetFirstName())
	h.router.NewRoute("/setUserName", false, h.handleSetUserName())
	h.logger.Debugf("Configuring callback commands router done")
}

func (h *callBackQueryHandler) handleMessage(u *tgbotapi.Update) {
	h.logger.Infof("Received CallBackData: dat")
}

func (h *callBackQueryHandler) handleCommand(u *tgbotapi.Update) {
	h.logger.Infof("Received CallBack Command: command=\"%s\" chatId=\"%d\"", u.CallbackQuery.Data, u.CallbackQuery.Message.Chat.ID)

	chatID := u.CallbackQuery.Message.Chat.ID
	user := h.store.User().FindUser(int(chatID))
	if user == nil {
		errMsg := tgbotapi.NewMessage(chatID, "Ошибка! Неизвестная команда")
		h.bot.Send(errMsg)
		return
	}
	if handler := h.router.GetHandler(u.CallbackQuery.Data); handler != nil {
		handler(user, u)
	} else {
		errMsg := tgbotapi.NewMessage(chatID, "Ошибка! Неизвестная команда")
		h.bot.Send(errMsg)
	}
}

func (h *callBackQueryHandler) updateIsCommand(u *tgbotapi.Update) bool {
	return strings.HasPrefix(u.CallbackQuery.Data, "/")
}

func (h *callBackQueryHandler) handleRegisterUser() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'RegisterUser'")

	type request struct {
		FirstName string `json:"firstName"`
		UserName  string `json:"userName"`
		TgID      int64  `json:"tgId"`
		From      string `json:"from"`
	}

	return func(user *model.User, u *tgbotapi.Update) {
		if user.WelcomeMessage.MessageID != u.CallbackQuery.Message.MessageID {
			return
		}

		// User registration
		req := &request{}
		req.FirstName = user.FirstName
		req.UserName = user.UserName
		req.TgID = user.UserID()
		req.From = "telegram"

		b, _ := json.Marshal(req)

		resp, err := http.Post(fmt.Sprintf("http://%s:%s/users", h.config.QaskAddress, h.config.QaskPort), "application/json", bytes.NewReader(b))
		if err != nil {
			h.internalError(user.UserID(), err)
			return
		}

		if resp.StatusCode != http.StatusCreated {
			// Error while register
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// Internal telegram bot error
				h.internalError(user.UserID(), err)
				return
			}

			h.internalError(user.UserID(), errors.New(string(b)))
			return
		}

		user.Registered = true

		message := model.WelcomeMessageAfterRegister(user)
		user.WelcomeMessageHead = message
		h.bot.Send(message.Msg)

		registeredHelpText :=
			`Вам доступны следующие команды:
/play - играть
/profile - настройки профиля
/newpass - сгенерировать новый пароль
`

		msg := tgbotapi.NewMessage(user.UserID(), registeredHelpText)
		h.bot.Send(msg)
	}
}

func (h *callBackQueryHandler) handleProfile() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'Profile'")

	return func(user *model.User, u *tgbotapi.Update) {
		if user.WelcomeMessage.MessageID != u.CallbackQuery.Message.MessageID {
			return
		}

		message := model.ProfileMain(user)
		user.ProfileMessageHead = message
		user.ProfileMessage, _ = h.bot.Send(message.Msg)
	}
}

func (h *callBackQueryHandler) handleGetQuestion() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'GetQuestion'")

	return func(user *model.User, u *tgbotapi.Update) {
		question := model.GetQuestion(user, h.config.QaskAddress, h.config.QaskPort)
		if question == nil {
			return
		}

		user.Question = question
		msg := tgbotapi.NewMessage(user.UserID(), user.Question.Question)

		var rows = make([][]tgbotapi.InlineKeyboardButton, 0)

		btnAnswer := makeButton("/showAnswer", "Показать ответ")
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnAnswer))

		btnReport := makeButton("/sendReport", "Сообщить о проблеме")
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnReport))

		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

		user.QuestionMessage, _ = h.bot.Send(msg)
	}
}

func (h *callBackQueryHandler) handleShowAnswer() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'ShowAnswer'")
	return func(user *model.User, u *tgbotapi.Update) {
		if user.QuestionMessage.MessageID == 0 {
			return
		}

		msg := tgbotapi.NewEditMessageText(u.CallbackQuery.Message.Chat.ID, user.QuestionMessage.MessageID, user.Question.Answer)

		var rows = make([][]tgbotapi.InlineKeyboardButton, 0)

		btnQuestion := makeButton("/showQuestion", "Показать вопрос")
		if user.Question.Comment != "" {
			btnComment := makeButton("/showComment", "Показать комментарий")
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnQuestion, btnComment))
		} else {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnQuestion))
		}

		btnReport := makeButton("/sendReport", "Сообщить о проблеме")
		btnGetQuestion := makeButton("/getQuestion", "Следующий вопрос")
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnReport, btnGetQuestion))

		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg.ReplyMarkup = &replyMarkup

		h.bot.Send(msg)
	}
}

func (h *callBackQueryHandler) handleShowQuestion() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'ShowQuestion'")
	return func(user *model.User, u *tgbotapi.Update) {
		msg := tgbotapi.NewEditMessageText(u.CallbackQuery.Message.Chat.ID, user.QuestionMessage.MessageID, user.Question.Question)

		var rows = make([][]tgbotapi.InlineKeyboardButton, 0)

		btnAnswer := makeButton("/showAnswer", "Показать ответ")
		if user.Question.Comment != "" {
			btnComment := makeButton("/showComment", "Показать комментарий")
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnAnswer, btnComment))
		} else {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnAnswer))
		}

		btnReport := makeButton("/sendReport", "Сообщить о проблеме")
		btnGetQuestion := makeButton("/getQuestion", "Следующий вопрос")
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnReport, btnGetQuestion))

		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg.ReplyMarkup = &replyMarkup

		h.bot.Send(msg)
	}
}

func (h *callBackQueryHandler) handleShowComment() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'ShowComment'")
	return func(user *model.User, u *tgbotapi.Update) {
		msg := tgbotapi.NewEditMessageText(u.CallbackQuery.Message.Chat.ID, user.QuestionMessage.MessageID, user.Question.Comment)

		var rows = make([][]tgbotapi.InlineKeyboardButton, 0)

		btnAnswer := makeButton("/showAnswer", "Показать ответ")
		if user.Question.Comment != "" {
			btnComment := makeButton("/showComment", "Показать комментарий")
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnAnswer, btnComment))
		} else {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnAnswer))
		}

		btnReport := makeButton("/sendReport", "Сообщить о проблеме")
		btnGetQuestion := makeButton("/getQuestion", "Следующий вопрос")
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnReport, btnGetQuestion))

		replyMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg.ReplyMarkup = &replyMarkup

		h.bot.Send(msg)

	}
}

func (h *callBackQueryHandler) handleGetMathProblem() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'GetMathProblem'")
	return func(user *model.User, u *tgbotapi.Update) {
		h.unavailableCommand(user.UserID())
	}
}

func (h *callBackQueryHandler) handleSettings() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'Settings'")
	return func(user *model.User, u *tgbotapi.Update) {
		message := model.GameMainSettingsMessage(user)
		message.Prev = user.PlayMessageHead
		user.PlayMessageHead = message
		h.bot.Send(message.Msg)
	}
}

func (h *callBackQueryHandler) handleBack() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'Back'")
	return func(user *model.User, u *tgbotapi.Update) {
		message := user.PlayMessageHead.Prev
		user.PlayMessageHead = message
		h.bot.Send(message.Msg)
	}
}

func (h *callBackQueryHandler) handleSubscriptions() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'Subscriptions'")
	return func(user *model.User, u *tgbotapi.Update) {
		message := model.GameSubscriptionsSettingsMessage(user)
		message.Prev = user.PlayMessageHead
		user.PlayMessageHead = message
		h.bot.Send(message.Msg)
	}
}

func (h *callBackQueryHandler) handleSendReport() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'SendReport'")

	return func(user *model.User, u *tgbotapi.Update) {
		msg := tgbotapi.NewMessage(user.UserID(), "Окей, введите сообщение о проблеме")

		user.WriteTo = &user.ReportMessage

		h.bot.Send(msg)
	}
}

func (h *callBackQueryHandler) handleSetFirstName() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'SetFirstName'")

	return func(user *model.User, u *tgbotapi.Update) {
		msg := tgbotapi.NewMessage(user.UserID(), "Окей, введите новое имя")

		user.WriteTo = &user.FirstName

		h.bot.Send(msg)
	}
}

func (h *callBackQueryHandler) handleSetUserName() router.RouterHandler {
	h.logger.Debugf("Register callback handler 'SetUserName'")

	return func(user *model.User, u *tgbotapi.Update) {
		msg := tgbotapi.NewMessage(user.UserID(), "Окей, введите новое имя пользователя")

		user.WriteTo = &user.UserName

		h.bot.Send(msg)
	}
}

/*
func (b *tgbot) handleCallbackQuery(update *tgbotapi.Update) {
	if update.CallbackQuery.Data == "/hidden_answer" {
		mssageConfig := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, lastMessage.MessageID, model.TestQuestion().Question)
		var rows = make([][]tgbotapi.InlineKeyboardButton, 0)
		strAnswer := "/show_answer"
		btnAnswer := tgbotapi.InlineKeyboardButton{
			Text:         "Показать ответ",
			CallbackData: &strAnswer,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnAnswer))

		strReport := "/send_report"
		btnReport := tgbotapi.InlineKeyboardButton{
			Text:         "Сообщить о проблеме",
			CallbackData: &strReport,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnReport))

		newReplyMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg.ReplyMarkup = &newReplyMarkup
		b.bot.Send(msg)

	} else {
		msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, lastMessage.MessageID, model.TestQuestion().Answer)
		var rows = make([][]tgbotapi.InlineKeyboardButton, 0)
		strAnswer := "/hidden_answer"
		btnAnswer := tgbotapi.InlineKeyboardButton{
			Text:         "Скрыть ответ",
			CallbackData: &strAnswer,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnAnswer))

		strReport := "/send_report"
		btnReport := tgbotapi.InlineKeyboardButton{
			Text:         "Сообщить о проблеме",
			CallbackData: &strReport,
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btnReport))

		newReplyMarkup := tgbotapi.NewInlineKeyboardMarkup(rows...)
		msg.ReplyMarkup = &newReplyMarkup
		b.bot.Send(msg)

	}

}
*/

func (h *callBackQueryHandler) unavailableCommand(chatID int64) {
	err := errors.New("Недоступная команда")
	h.internalError(chatID, err)
}

func (h *callBackQueryHandler) internalError(chatID int64, err error) {
	errorMessage := fmt.Sprintf("Произошла внутренняя ошибка:\n\"%s\"\nПожалуйста, повторите попытку позже.", err)
	msg := tgbotapi.NewMessage(chatID, errorMessage)
	h.bot.Send(msg)
}

func makeButton(callbackData string, label string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{
		Text:         label,
		CallbackData: &callbackData,
	}
}

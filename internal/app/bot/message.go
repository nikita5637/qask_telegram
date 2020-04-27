package bot

import (
	"qask_telegram/internal/app/model"
	"qask_telegram/internal/app/router"
	"qask_telegram/internal/app/store"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type messageHandler struct {
	bot    *tgbotapi.BotAPI
	logger *logrus.Logger
	router *router.Router
	store  store.Store
}

func newMessageHandler(bot *tgbotapi.BotAPI, logger *logrus.Logger, store store.Store) *messageHandler {
	mH := &messageHandler{
		bot:    bot,
		logger: logger,
		router: router.NewRouter(logger),
		store:  store,
	}

	mH.configureRouter()

	return mH
}

func (h *messageHandler) configureRouter() {
	h.logger.Debugf("Configuring message commands router ...")
	// Registering new routes (path, isPublic, handler)
	h.router.NewRoute("/help", true, h.handleHelp())
	h.router.NewRoute("/start", true, h.handleStart())
	h.router.NewRoute("/play", true, h.handlePlay())
	h.router.NewRoute("/report", true, h.handleReport())
	h.router.NewRoute("/profile", true, h.handleProfile())
	h.logger.Debugf("Configuring message commands router done")
}

func (h *messageHandler) handleMessage(u *tgbotapi.Update) {
	h.logger.Infof("Received Message: text=\"%s\" chatId=\"%d\"", u.Message.Text, u.Message.Chat.ID)

	chatID := u.Message.Chat.ID
	user := h.store.User().FindUser(int(chatID))

	if user.WriteTo != nil {
		*user.WriteTo = u.Message.Text
		user.WriteTo = nil
	}
}

func (h *messageHandler) handleCommand(u *tgbotapi.Update) {
	text := u.Message.Text
	chatID := u.Message.Chat.ID
	h.logger.Infof("Received Message Command: command=\"%s\" chatId=\"%d\"", text, chatID)

	user := h.store.User().FindUser(int(chatID))

	if user == nil {
		if u.Message.Text != "/start" && u.Message.Text != "/help" {
			h.unavailableCommand(chatID)
			return
		}
	}

	if handler := h.router.GetHandler(u.Message.Text); handler != nil {
		if h.router.CommandIsPublic(u.Message.Text) {
			handler(user, u)
		} else {
			h.unavailableCommand(chatID)
		}
	} else {
		h.unavailableCommand(chatID)
	}
}

func (h *messageHandler) updateIsCommand(u *tgbotapi.Update) bool {
	return strings.HasPrefix(u.Message.Text, "/")
}

func (h *messageHandler) handleHelp() router.RouterHandler {
	h.logger.Debugf("Register message handler 'Help'")

	unregisteredHelpText :=
		`Вам доступны следующие команды:
/help - справка
/start - регистрация
`
	unregisteredHelpMessage := tgbotapi.NewMessage(0, unregisteredHelpText)

	registeredHelpText :=
		`Вам доступны следующие команды:
/play - играть
/profile - настройки профиля
/newpass - сгенерировать новый пароль
`
	registeredHelpMessage := tgbotapi.NewMessage(0, registeredHelpText)

	return func(user *model.User, u *tgbotapi.Update) {
		chatID := u.Message.Chat.ID
		msg := unregisteredHelpMessage
		msg.ChatID = chatID

		if user == nil {
			msg = unregisteredHelpMessage
			msg.ChatID = chatID
		} else {
			if user.IsRegistered() == true {
				msg = registeredHelpMessage
				msg.ChatID = int64(user.UserId)
			}
		}

		h.bot.Send(msg)
	}
}

func (h *messageHandler) handleStart() router.RouterHandler {
	h.logger.Debugf("Register message handler 'Start'")

	return func(user *model.User, u *tgbotapi.Update) {
		if user == nil {
			h.logger.Infof("Register new user: ID=\"%d\" FirstName=\"%s\" LastName=\"%s\" UserName \"%s\" LanguageCode=\"%s\" IsBot=\"%t\"",
				u.Message.From.ID,
				u.Message.From.FirstName,
				u.Message.From.LastName,
				u.Message.From.UserName,
				u.Message.From.LanguageCode,
				u.Message.From.IsBot)

			user = h.store.User().CreateUser(u.Message.From.ID)
			user.FirstName = u.Message.From.FirstName
			user.UserName = u.Message.From.UserName

			message := model.WelcomeMessage(user)
			user.WelcomeMessageHead = message
			user.WelcomeMessage, _ = h.bot.Send(message.Msg)
		} else {
			if user.Registered {
				msg := tgbotapi.NewMessage(user.UserID(), "Вы уже зарегистрированы!")
				h.bot.Send(msg)
			} else {
				message := model.WelcomeMessage(user)
				user.WelcomeMessageHead = message
				user.WelcomeMessage, _ = h.bot.Send(message.Msg)
			}
		}
	}
}

func (h *messageHandler) handlePlay() router.RouterHandler {
	h.logger.Debugf("Register handler 'Play'")

	return func(user *model.User, u *tgbotapi.Update) {
		//chatId := user.UserID()
		if !user.Registered {
			h.unavailableCommand(user.UserID())
			return
		}

		/*
			var keyboardMarkup = make([][]tgbotapi.InlineKeyboardButton, 0)

			if user.QuestSubscribtion == true {
				btnGetQuestion := tgbotapi.NewInlineKeyboardButtonData("Случайный вопрос", "/getQuestion")
				keyboardMarkup = append(keyboardMarkup, tgbotapi.NewInlineKeyboardRow(btnGetQuestion))
			}

			if user.MathProblemSubscribtion == true {
				btnGetMathProblem := tgbotapi.NewInlineKeyboardButtonData("Математическая задача", "/getMathProblem")
				keyboardMarkup = append(keyboardMarkup, tgbotapi.NewInlineKeyboardRow(btnGetMathProblem))
			}

			btnSettings := tgbotapi.NewInlineKeyboardButtonData("Настройки игры", "/settings")
			keyboardMarkup = append(keyboardMarkup, tgbotapi.NewInlineKeyboardRow(btnSettings))

			msg := tgbotapi.NewMessage(chatId, "Выберите действие:")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboardMarkup...)
		*/

		message := model.PlayMessage(user)
		user.PlayMessageHead = message
		user.PlayMessage, _ = h.bot.Send(message.Msg)
	}
}

func (h *messageHandler) handleReport() router.RouterHandler {
	h.logger.Debugf("Register handler 'Report'")

	return func(user *model.User, u *tgbotapi.Update) {
		h.unavailableCommand(user.UserID())
	}
}

func (h *messageHandler) handleProfile() router.RouterHandler {
	h.logger.Debugf("Register message handler 'Profile'")

	return func(user *model.User, u *tgbotapi.Update) {
		message := model.ProfileMain(user)
		user.ProfileMessageHead = message
		user.ProfileMessage, _ = h.bot.Send(message.Msg)
	}
}

func (h *messageHandler) unavailableCommand(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Недоступная команда")
	h.bot.Send(msg)
}

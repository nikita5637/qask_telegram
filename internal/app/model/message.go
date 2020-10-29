package model

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Message is a structure, that allow to use a multi-level editable message
type Message struct {
	Msg  tgbotapi.Chattable
	Prev *Message
}

//PlayMessage is a first-level game editable message
func PlayMessage(user *User) *Message {
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

	msg := tgbotapi.NewMessage(user.UserID(), "Выберите действие:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboardMarkup...)

	newMessage := &Message{
		Msg:  &msg,
		Prev: nil,
	}

	return newMessage
}

//GameMainSettingsMessage ...
func GameMainSettingsMessage(user *User) *Message {
	msg := tgbotapi.NewEditMessageText(user.UserID(), user.PlayMessage.MessageID, "Настройки игры")

	btn1 := tgbotapi.NewInlineKeyboardButtonData("Подписки", "/subscribtions")
	btn1Row := tgbotapi.NewInlineKeyboardRow(btn1)

	btnBack := tgbotapi.NewInlineKeyboardButtonData("<< назад", "/back")
	btnBackRow := tgbotapi.NewInlineKeyboardRow(btnBack)

	rows := tgbotapi.NewInlineKeyboardMarkup(btn1Row, btnBackRow)

	msg.ReplyMarkup = &rows

	return &Message{
		Msg: &msg,
	}
}

//GameSubscriptionsSettingsMessage ...
func GameSubscriptionsSettingsMessage(user *User) *Message {
	msg := tgbotapi.NewEditMessageText(user.UserID(), user.PlayMessage.MessageID, "Настройки подписок")

	btn1 := tgbotapi.NewInlineKeyboardButtonData("Получать вопросы", "/subscribeQuestions")
	btn1Row := tgbotapi.NewInlineKeyboardRow(btn1)

	btn2 := tgbotapi.NewInlineKeyboardButtonData("Получать математические задачи", "/subscribeMathProblems")
	btn2Row := tgbotapi.NewInlineKeyboardRow(btn2)

	btnBack := tgbotapi.NewInlineKeyboardButtonData("<< назад", "/back")
	btnBackRow := tgbotapi.NewInlineKeyboardRow(btnBack)

	rows := tgbotapi.NewInlineKeyboardMarkup(btn1Row, btn2Row, btnBackRow)

	msg.ReplyMarkup = &rows

	return &Message{
		Msg: &msg,
	}
}

// ProfileMain ...
func ProfileMain(user *User) *Message {
	msgProfile := fmt.Sprintf("Настройки профиля")

	strSetFirstName := fmt.Sprintf("Имя [%s]", user.FirstName)
	btnSetFirstName := tgbotapi.NewInlineKeyboardButtonData(strSetFirstName, "/setFirstName")
	btnSetFirstNameRow := tgbotapi.NewInlineKeyboardRow(btnSetFirstName)

	strSetUserName := fmt.Sprintf("Имя пользователя [%s]", user.UserName)
	btnSetUserName := tgbotapi.NewInlineKeyboardButtonData(strSetUserName, "/setUserName")
	btnSetUserNameRow := tgbotapi.NewInlineKeyboardRow(btnSetUserName)

	msgProfileKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(btnSetFirstNameRow, btnSetUserNameRow)

	msg := tgbotapi.NewMessage(int64(user.UserId), msgProfile)
	msg.ReplyMarkup = msgProfileKeyboardMarkup

	return &Message{
		Msg:  &msg,
		Prev: nil,
	}
}

//WelcomeMessage is a "start" message a user recieves when sending message "/start"
func WelcomeMessage(user *User) *Message {
	msgWelcome := fmt.Sprintf(
		`Добро пожаловать, %s!
Для игры необходимо зарегистрироваться. Попробуй сделать это прямо сейчас, нажав кнопку "Зарегистрироваться".`, user.FirstName)

	btnRegister := tgbotapi.NewInlineKeyboardButtonData("Зарегистрироваться", "/register")
	btnRegisterRow := tgbotapi.NewInlineKeyboardRow(btnRegister)

	btnProfileSettings := tgbotapi.NewInlineKeyboardButtonData("Настройки профиля", "/profile")
	btnProfileSettingsRow := tgbotapi.NewInlineKeyboardRow(btnProfileSettings)

	msgWelcomeKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(btnRegisterRow, btnProfileSettingsRow)

	msg := tgbotapi.NewMessage(int64(user.UserId), msgWelcome)
	msg.ReplyMarkup = msgWelcomeKeyboardMarkup

	return &Message{
		Msg:  &msg,
		Prev: nil,
	}
}

//WelcomeMessageAfterRegister ...
func WelcomeMessageAfterRegister(user *User) *Message {
	msg := tgbotapi.NewEditMessageText(user.UserID(), user.WelcomeMessage.MessageID, "Регистрация прошла успешно")

	btn1 := tgbotapi.NewInlineKeyboardButtonData("Настройки профиля", "/profile")
	btn1Row := tgbotapi.NewInlineKeyboardRow(btn1)

	rows := tgbotapi.NewInlineKeyboardMarkup(btn1Row)

	msg.ReplyMarkup = &rows

	return &Message{
		Msg: &msg,
	}
}

//WelcomeMessageAfterSendingReport ...
func WelcomeMessageAfterSendingReport(user *User) *Message {
	msg := tgbotapi.NewMessage(user.UserID(), "Ваше обращение получено. Спасибо!")

	return &Message{
		Msg: &msg,
	}
}

package bot

import (
	"qask_telegram/internal/app/store"
	"qask_telegram/internal/app/store/cache"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type updateHandler interface {
	handleMessage(*tgbotapi.Update)
	handleCommand(*tgbotapi.Update)
	updateIsCommand(*tgbotapi.Update) bool
}

type tgbot struct {
	bot                  *tgbotapi.BotAPI
	logger               *logrus.Logger
	updChan              *tgbotapi.UpdatesChannel
	store                store.Store
	callBackQueryHandler *callBackQueryHandler
	messageHandler       *messageHandler
}

//Start ...
func Start(token string) error {
	bot, err := startBot(token)
	if err != nil {
		return err
	}

	logger := logrus.New()
	level, err := logrus.ParseLevel("debug")
	if err != nil {
		return err
	}

	logger.SetLevel(level)

	st := cache.New(logger)

	bot.logger = logger
	bot.store = st

	bot.callBackQueryHandler = newCallBackQueryHandler(bot.bot, logger, st)
	bot.messageHandler = newMessageHandler(bot.bot, logger, st)

	for update := range *bot.updChan {
		if update.CallbackQuery == nil && update.Message == nil && update.ChannelPost == nil {
			continue
		} else if update.CallbackQuery != nil {
			go bot.ServeUpdate(&update, bot.callBackQueryHandler)
		} else if update.Message != nil {
			go bot.ServeUpdate(&update, bot.messageHandler)
		}
	}
	return nil
}

func startBot(token string) (*tgbot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	uc := tgbotapi.NewUpdate(0)
	uc.Timeout = 60

	updatesChan, err := bot.GetUpdatesChan(uc)
	if err != nil {
		return nil, err
	}

	return &tgbot{
		bot:     bot,
		logger:  nil,
		updChan: &updatesChan,
	}, nil
}

func (b *tgbot) ServeUpdate(update *tgbotapi.Update, handler updateHandler) {
	if handler.updateIsCommand(update) {
		handler.handleCommand(update)
	} else {
		handler.handleMessage(update)
	}
	/*
		if update.CallbackQuery != nil {
			b.logger.Infof("Received callback query: data=\"%s\" chatid=\"%d\" msgid=\"%d\"", update.CallbackQuery.Data, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
			b.PrintSenderData(update.Message.From)

			//chatId := update.CallbackQuery.Message.Chat.ID
			//user := b.store.User().FindUser(chatId)
			b.handleCallBackQuery(update)
		} else if update.Message != nil {
			b.logger.Infof("Received text message: text=\"%s\" chatid=\"%d\"", update.Message.Text, update.Message.Chat.ID)
			b.PrintSenderData(update.Message.From)
			chatId := update.Message.Chat.ID
			user := b.store.User().FindUser(chatId)
			if user != nil {
				b.logger.Debugf("\tUsername=\"%s\"", user.Name)
			} else {
				user = b.store.User().CreateUser(chatId)
			}

			if handler := b.router.GetHandler(update.Message.Text); handler != nil {
				handler(user, update)
			} else {
				b.logger.Infof("Can not find handler for message: text=\"%s\" chatid\"%s\"", update.Message.Text, update.Message.Chat.ID)
			}
		}
	*/
}

/*

func (b *tgbot) SendError(chatId int64, errString string) {
	msg := tgbotapi.NewMessage(chatId, errString)
	b.bot.Send(msg)
}

func (b *tgbot) PrintSenderData(user *tgbotapi.User) {
	b.logger.Infof("Sender data: ID=\"%d\" FirstName=\"%s\" LastName=\"%s\" UserName=\"%s\" LanguageCode=\"%s\" IsBot=\"%t\"",
		user.ID, user.FirstName, user.LastName, user.UserName, user.LanguageCode, user.IsBot)
}
*/

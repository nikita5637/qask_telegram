package model

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type userPublic struct {
	FirstName string `json:"firstName"`
	UserName  string `json:"userName"`
}

type userPrivate struct {
	DBID                    int  `json:"dbId"`
	UserId                  int  `json:"UserId"`
	Registered              bool `json:"registered"`
	State                   int  `json:"state"`
	QuestSubscribtion       bool
	MathProblemSubscribtion bool
	WelcomeMessage          tgbotapi.Message
	WelcomeMessageHead      *Message
	ProfileMessage          tgbotapi.Message
	ProfileMessageHead      *Message
	PlayMessage             tgbotapi.Message
	PlayMessageHead         *Message
	ReportMessage           string
	QuestionMessage         tgbotapi.Message
	Question                *Question
	WriteTo                 *string
}

type User struct {
	userPublic
	userPrivate
}

func (u *User) Validate() error {
	if u.FirstName == "" {
		return errors.New("Username is empty")
	}
	return nil
}

func (u *User) IsRegistered() bool {
	if u == nil {
		return false
	}

	return u.Registered == true
}

func (u *User) UserID() int64 {
	return int64(u.UserId)
}

package cache

import (
	"github.com/sirupsen/logrus"
	"qask_telegram/internal/app/model"
)

type UserRepository struct {
	users  map[int]*model.User
	logger *logrus.Logger
}

func (u *UserRepository) CreateUser(id int) *model.User {

	if user := u.FindUser(id); user != nil {
		return user
	}

	newUser := &model.User{}

	newUser.UserId = id
	newUser.Registered = false
	newUser.QuestSubscribtion = true
	newUser.MathProblemSubscribtion = true

	u.users[id] = newUser

	return newUser
}

func (u *UserRepository) RegisterUser(chatid int, name string) error {
	return nil
}

/*
func (u *UserRepository) RegisterUser(chatid int, name string) error {
	user := u.FindUser(chatid)
	if user == nil {
		return errors.New("User not found")
	}

	if user.Registered == true {
		return errors.New("User already registered")
	}

	u.logger.Infof("Registering new user with chatid=\"%d\", name=\"%s\" ...", chatid, name)

	user.FirstName = name
	if err := user.Validate(); err != nil {
		return err
	}
	user.Registered = true

	u.logger.Infof("Registering new user with chatid=\"%d\", name=\"%s\" done", chatid, name)
	return nil
}
*/

func (u *UserRepository) FindUser(chatid int) *model.User {
	if _, ok := u.users[chatid]; ok {
		u.logger.Debugf("User with chat id '%d' found", chatid)
	} else {
		u.logger.Debugf("User with chat id '%d' not found", chatid)
	}
	return u.users[chatid]
}

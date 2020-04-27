package store

import (
	"qask_telegram/internal/app/model"
)

type UserRepository interface {
	CreateUser(int) *model.User
	RegisterUser(int, string) error
	FindUser(int) *model.User
}

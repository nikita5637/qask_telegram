package cache

import (
	"github.com/sirupsen/logrus"
	"qask_telegram/internal/app/model"
	"qask_telegram/internal/app/store"
)

type Store struct {
	userRepository *UserRepository
	logger         *logrus.Logger
}

func New(logger *logrus.Logger) *Store {
	return &Store{
		logger: logger,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		users:  make(map[int]*model.User),
		logger: s.logger,
	}

	return s.userRepository
}

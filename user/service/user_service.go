package service

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-els-tg-bot/models"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type userService struct {
	userRepo user.UserRepository
	logger   *logrus.Logger
}

// NewUserService creates a new instance of user service
func NewUserService(userRepo user.UserRepository, logger *logrus.Logger) user.UserService {
	return &userService{
		userRepo,
		logger,
	}
}

// Create creates a new user
func (s *userService) Create(chatID int64) (*models.User, error) {
	user := models.User{
		ChatID: chatID,
	}
	id, err := s.userRepo.Store(&user)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	user.ID = id

	userJSON, _ := json.Marshal(user)
	s.logger.Info("User created: " + string(userJSON))

	return &user, nil
}

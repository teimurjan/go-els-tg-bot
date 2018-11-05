package service

import (
	"github.com/teimurjan/go-els-tg-bot/models"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type userService struct {
	userRepo user.UserRepository
}

func NewUserService(userRepo user.UserRepository) *userService {
	return &userService{
		userRepo,
	}
}

func (us *userService) Create(chatID int64) (*models.User, error) {
	user := models.User{
		ChatID: chatID,
	}
	id, err := us.userRepo.Store(&user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return &user, nil
}

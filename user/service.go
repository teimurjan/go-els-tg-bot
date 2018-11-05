package user

import "github.com/teimurjan/go-els-tg-bot/models"

type UserService interface {
	Create(chatID int64) (*models.User, error)
}

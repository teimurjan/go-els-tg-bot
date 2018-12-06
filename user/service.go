package user

import "github.com/teimurjan/go-els-tg-bot/models"

// UserService is an interface for user service
type UserService interface {
	Create(chatID int64) (*models.User, error)
}

package user

import (
	"github.com/teimurjan/go-els-tg-bot/models"
)

type UserRepository interface {
	GetByID(id int64) (*models.User, error)
	GetByChatID(chatID int64) (*models.User, error)
	Store(u *models.User) (int64, error)
	GetAll() ([]*models.User, error)
}

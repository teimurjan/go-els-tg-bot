package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/teimurjan/go-els-tg-bot/models"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type mysqlUserRepository struct {
	conn *sqlx.DB
}

// NewMysqlUserRepository creates a new instance of mysql repository for user
func NewMysqlUserRepository(conn *sqlx.DB) user.UserRepository {
	return &mysqlUserRepository{conn}
}

// GetByID returns a user by id
func (m *mysqlUserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := m.conn.Get(&user, `
		SELECT *
		FROM users
		WHERE id=?
	`, id)

	return &user, err
}

// GetByChatID returns a user by chat_id
func (m *mysqlUserRepository) GetByChatID(chatID int64) (*models.User, error) {
	var user models.User
	err := m.conn.Get(&user, `
		SELECT *
		FROM users
		WHERE chat_id=?
	`, chatID)
	return &user, err
}

// Store adds a new user
func (m *mysqlUserRepository) Store(user *models.User) (int64, error) {
	currentTime := time.Now().UTC()

	res, err := m.conn.Exec(`
		INSERT INTO users
		(chat_id, created, modified)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE modified=?;
	`, user.ChatID, currentTime, currentTime, currentTime)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}

// Update updates a user
func (m *mysqlUserRepository) Update(user *models.User) error {
	currentTime := time.Now().UTC()

	_, err := m.conn.Exec(`
		UPDATE users
		SET language=?, modified=?
		WHERE id=?;
	`, user.Language, currentTime, user.ID)

	return err
}

// GetAll returns all existing users
func (m *mysqlUserRepository) GetAll() ([]*models.User, error) {
	var users []*models.User

	err := m.conn.Select(&users, `
		SELECT * 
		FROM users;
	`)

	if err != nil {
		return nil, err
	}

	return users, nil
}

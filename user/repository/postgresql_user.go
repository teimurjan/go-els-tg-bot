package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/teimurjan/go-els-tg-bot/models"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type postgresqlUserRepository struct {
	conn *sqlx.DB
}

// NewPostgresqlUserRepository creates a new instance of postgresql repository for user
func NewPostgresqlUserRepository(conn *sqlx.DB) user.UserRepository {
	return &postgresqlUserRepository{conn}
}

// GetByID returns a user by id
func (m *postgresqlUserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := m.conn.Get(&user, `
		SELECT *
		FROM users
		WHERE id=?;
	`, id)

	return &user, err
}

// GetByChatID returns a user by chat_id
func (m *postgresqlUserRepository) GetByChatID(chatID int64) (*models.User, error) {
	var user models.User
	err := m.conn.Get(&user, `
		SELECT *
		FROM users
		WHERE chat_id=$1;
	`, chatID)
	return &user, err
}

// Store adds a new user
func (m *postgresqlUserRepository) Store(user *models.User) (int64, error) {
	currentTime := time.Now().UTC()

	var id int64
	err := m.conn.QueryRow(`
		INSERT INTO users
		(chat_id, created, modified)
		VALUES ($1, $2, $2)
		ON CONFLICT (chat_id) DO UPDATE 
		SET modified=$2
		RETURNING id;
	`, user.ChatID, currentTime).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

// Update updates a user
func (m *postgresqlUserRepository) Update(user *models.User) error {
	currentTime := time.Now().UTC()

	_, err := m.conn.Exec(`
		UPDATE users
		SET language=$1, modified=$2
		WHERE id=$3;
	`, user.Language, currentTime, user.ID)

	return err
}

// GetAll returns all existing users
func (m *postgresqlUserRepository) GetAll() ([]*models.User, error) {
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

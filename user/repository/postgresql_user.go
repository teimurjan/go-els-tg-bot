package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/teimurjan/go-els-tg-bot/models"
)

type postgresqlUserRepository struct {
	conn *sqlx.DB
}

func NewPostgresqlUserRepository(conn *sqlx.DB) *postgresqlUserRepository {
	return &postgresqlUserRepository{conn}
}

func (m *postgresqlUserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := m.conn.Get(&user, `
		SELECT *
		FROM users
		WHERE id=?;
	`, id)

	return &user, err
}

func (m *postgresqlUserRepository) GetByChatID(chatID int64) (*models.User, error) {
	var user models.User
	err := m.conn.Get(&user, `
		SELECT *
		FROM users
		WHERE chat_id=$1;
	`, chatID)
	return &user, err
}

func (m *postgresqlUserRepository) Store(user *models.User) (int64, error) {
	currentTime := time.Now().UTC()

	res, err := m.conn.Exec(`
		INSERT INTO users
		(chat_id, created, modified)
		VALUES ($1, $2, $2)
		ON CONFLICT (chat_id) DO UPDATE 
  		SET modified=$2;
	`, user.ChatID, currentTime)

	if err != nil {
		return 0, err
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return id, err
}

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

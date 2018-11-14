package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/teimurjan/go-els-tg-bot/models"
)

type postgresqlTrackingRepository struct {
	conn *sqlx.DB
}

func NewPostgresqlTrackingRepository(conn *sqlx.DB) *postgresqlTrackingRepository {
	return &postgresqlTrackingRepository{conn}
}

func (m *postgresqlTrackingRepository) GetByID(id int64) (*models.Tracking, error) {
	var tracking models.Tracking
	err := m.conn.Get(&tracking, `
		SELECT *
		FROM trackings 
		WHERE id=$1;
	`, id)

	return &tracking, err
}

func (m *postgresqlTrackingRepository) GetForUser(userID int64) ([]*models.Tracking, error) {
	var trackings []*models.Tracking
	err := m.conn.Select(&trackings, `
		SELECT *
		FROM trackings
		WHERE user_id=$1;
	`, userID)
	return trackings, err
}

func (m *postgresqlTrackingRepository) Store(t *models.Tracking) (int64, error) {
	currentTime := time.Now().UTC()

	var id int64
	err := m.conn.QueryRow(`
		INSERT INTO trackings
		(name, value, status, user_id, created, modified)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id;
	`, t.Name, t.Value, t.Status, t.UserID, currentTime).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

func (m *postgresqlTrackingRepository) Update(t *models.Tracking) error {
	currentTime := time.Now().UTC()

	_, err := m.conn.Exec(`
		UPDATE trackings
		SET name=$1, user_id=$2, status=$3, modified=$4
		WHERE value=$5;
	`, t.Name, t.UserID, t.Status, currentTime, t.Value)

	return err
}

func (m *postgresqlTrackingRepository) Delete(ID int64) error {
	_, err := m.conn.Exec(`
		DELETE FROM trackings
		WHERE id=$1;
	`, ID)

	return err
}

package repository

import (
	"time"

	"github.com/teimurjan/go-els-tg-bot/tracking"

	"github.com/jmoiron/sqlx"
	"github.com/teimurjan/go-els-tg-bot/models"
)

type mysqlTrackingRepository struct {
	conn *sqlx.DB
}

// NewMysqlTrackingRepository creates new mysqlTrackingRepository instance
func NewMysqlTrackingRepository(conn *sqlx.DB) tracking.TrackingRepository {
	return &mysqlTrackingRepository{conn}
}

func (m *mysqlTrackingRepository) GetByID(id int64) (*models.Tracking, error) {
	var tracking models.Tracking
	err := m.conn.Get(&tracking, `
		SELECT *
		FROM trackings 
		WHERE id=?
	`, id)

	return &tracking, err
}

func (m *mysqlTrackingRepository) GetForUser(userID int64) ([]*models.Tracking, error) {
	var trackings []*models.Tracking
	err := m.conn.Select(&trackings, `
		SELECT *
		FROM trackings
		WHERE user_id=?
	`, userID)

	return trackings, err
}

func (m *mysqlTrackingRepository) Store(t *models.Tracking) (int64, error) {
	currentTime := time.Now().UTC()

	res, err := m.conn.Exec(`
		INSERT INTO trackings
		(name, value, status, user_id, created, modified)
		VALUES (?, ?, ?, ?, ?, ?)
	`, t.Name, t.Value, t.Status, t.UserID, currentTime, currentTime)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}

func (m *mysqlTrackingRepository) Update(t *models.Tracking) error {
	currentTime := time.Now().UTC()

	_, err := m.conn.Exec(`
		UPDATE trackings SET
		name=?, status=?, user_id=?, modified=?
		WHERE value=?
	`, t.Name, t.Status, t.UserID, currentTime, t.Value)

	return err
}

func (m *mysqlTrackingRepository) Delete(ID int64) error {
	_, err := m.conn.Exec(`
		DELETE FROM trackings
		WHERE id=?
	`, ID)

	return err
}

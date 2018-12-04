package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/teimurjan/go-els-tg-bot/models"
)

type mysqlAddTrackingDialogRepository struct {
	conn *sqlx.DB
}

func NewMysqlAddTrackingDialogRepository(conn *sqlx.DB) *mysqlAddTrackingDialogRepository {
	return &mysqlAddTrackingDialogRepository{conn}
}

func (m *mysqlAddTrackingDialogRepository) GetForUser(userID int64) (*models.AddTrackingDialog, error) {
	var dialog models.AddTrackingDialog
	err := m.conn.Get(&dialog, `
		SELECT *
		FROM add_tracking_dialogs
		WHERE user_id=?
	`, userID)

	return &dialog, err
}

func (m *mysqlAddTrackingDialogRepository) Store(t *models.AddTrackingDialog) (int64, error) {
	currentTime := time.Now().UTC()

	res, err := m.conn.Exec(`
		INSERT INTO add_tracking_dialogs
		(step, user_id, created, modified)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE modified=?, step=?;
	`, t.Step, t.UserID, currentTime, currentTime, currentTime, t.Step)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}

func (m *mysqlAddTrackingDialogRepository) Update(t *models.AddTrackingDialog) error {
	currentTime := time.Now().UTC()

	_, err := m.conn.Exec(`
		UPDATE add_tracking_dialogs SET
		step=?, future_tracking_name=? modified=?
		WHERE id=?
	`, t.Step, t.FutureTrackingName, currentTime, t.ID)

	return err
}

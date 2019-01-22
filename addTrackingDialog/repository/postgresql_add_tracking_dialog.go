package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/teimurjan/go-els-tg-bot/addTrackingDialog"
	"github.com/teimurjan/go-els-tg-bot/models"
)

type postgresqlAddTrackingDialogRepository struct {
	conn *sqlx.DB
}

// NewPostgresqlAddTrackingDialogRepository creates new postgresqlAddTrackingDialogRepository instance
func NewPostgresqlAddTrackingDialogRepository(conn *sqlx.DB) addTrackingDialog.AddTrackingDialogRepository {
	return &postgresqlAddTrackingDialogRepository{conn}
}

func (m *postgresqlAddTrackingDialogRepository) GetForUser(userID int64) (*models.AddTrackingDialog, error) {
	var dialog models.AddTrackingDialog
	err := m.conn.Get(&dialog, `
		SELECT *
		FROM add_tracking_dialogs
		WHERE user_id=$1;
	`, userID)
	return &dialog, err
}

func (m *postgresqlAddTrackingDialogRepository) Store(d *models.AddTrackingDialog) (int64, error) {
	currentTime := time.Now().UTC()

	var id int64

	err := m.conn.QueryRow(`
		INSERT INTO add_tracking_dialogs
		(step, user_id, created, modified)
		VALUES ($1, $2, $3, $3)
		ON CONFLICT (user_id) DO UPDATE 
		SET step=$1, modified=$3
		RETURNING id;
	`, d.Step, d.UserID, currentTime).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

func (m *postgresqlAddTrackingDialogRepository) Update(d *models.AddTrackingDialog) error {
	currentTime := time.Now().UTC()

	_, err := m.conn.Exec(`
		UPDATE add_tracking_dialogs
		SET step=$1, future_tracking_name=$2, modified=$3
		WHERE id=$4;
	`, d.Step, d.FutureTrackingName, currentTime, d.ID)

	return err
}

package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/teimurjan/go-els-tg-bot/models"
	"github.com/teimurjan/go-els-tg-bot/usa-address"
)

type postgresqlUsaAddressRepository struct {
	conn *sqlx.DB
}

// NewPostgresqlUsaAddressRepository creates new postgresqlUsaAddressRepository instance
func NewPostgresqlUsaAddressRepository(conn *sqlx.DB) usaAddress.UsaAddressRepository {
	return &postgresqlUsaAddressRepository{conn}
}

func (m *postgresqlUsaAddressRepository) GetFirst() (*models.UsaAddress, error) {
	var usaAddress models.UsaAddress
	err := m.conn.Get(&usaAddress, `
		SELECT *
		FROM usaAddresses 
		LIMIT 1;
	`)

	return &usaAddress, err
}

func (m *postgresqlUsaAddressRepository) Store(a *models.UsaAddress) (int64, error) {
	currentTime := time.Now().UTC()

	var id int64
	err := m.conn.QueryRow(`
		INSERT INTO usaAddresses
		(street, city, state, zip, phone, created, modified)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id;
	`, a.Street, a.City, a.State, a.Zip, a.PhoneNumber, currentTime).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

func (m *postgresqlUsaAddressRepository) Update(a *models.UsaAddress) error {
	currentTime := time.Now().UTC()

	_, err := m.conn.Exec(`
		UPDATE usaAddresses
		SET street=$1, city=$2, state=$3, zip=$4, phone=$5, modified=$6;
	`, a.Street, a.City, a.State, a.Zip, a.PhoneNumber, currentTime)

	return err
}

func (m *postgresqlUsaAddressRepository) Delete() error {
	_, err := m.conn.Exec(`
		DELETE FROM usaAddresses;
	`)

	return err
}

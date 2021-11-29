package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/teimurjan/go-els-tg-bot/models"
	"github.com/teimurjan/go-els-tg-bot/usa-address"
)

type mysqlUsaAddressRepository struct {
	conn *sqlx.DB
}

// NewMysqlUsaAddressRepository creates new mysqlUsaAddressRepository instance
func NewMysqlUsaAddressRepository(conn *sqlx.DB) usaAddress.UsaAddressRepository {
	return &mysqlUsaAddressRepository{conn}
}

func (m *mysqlUsaAddressRepository) GetFirst() (*models.UsaAddress, error) {
	var usaAddress models.UsaAddress
	err := m.conn.Get(&usaAddress, `
		SELECT *
		FROM usaAddresses 
		LIMIT 1;
	`)

	return &usaAddress, err
}

func (m *mysqlUsaAddressRepository) Store(a *models.UsaAddress) (int64, error) {
	currentTime := time.Now().UTC()

	var id int64
	err := m.conn.QueryRow(`
		INSERT INTO usaAddresses
		(street, city, state, zip, phone, created, modified)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`, a.Street, a.City, a.State, a.Zip, a.PhoneNumber, currentTime, currentTime).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

func (m *mysqlUsaAddressRepository) Update(a *models.UsaAddress) error {
	currentTime := time.Now().UTC()

	_, err := m.conn.Exec(`
		UPDATE usaAddresses
		SET street=?, city=?, state=?, zip=?, phone=?, modified=?;
	`, a.Street, a.City, a.State, a.Zip, a.PhoneNumber, currentTime)

	return err
}

func (m *mysqlUsaAddressRepository) Delete() error {
	_, err := m.conn.Exec(`
		DELETE FROM usaAddresses;
	`)

	return err
}

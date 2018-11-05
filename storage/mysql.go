package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/teimurjan/go-els-tg-bot/config"
)

// NewMySQL creates a new MySQL instance
func NewMySQL(c *config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open("mysql", c.DatabaseURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

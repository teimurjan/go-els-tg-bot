package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/teimurjan/go-els-tg-bot/config"
)

// NewMySQL creates a new MySQL instance
func NewPostgreSQL(c *config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open("postgres", c.DatabaseURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

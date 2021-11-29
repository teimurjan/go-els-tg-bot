package models

import "time"

// UsaAddress is a model of ELS usa address
type UsaAddress struct {
	ID          int64     `db:"id" json:"id"`
	Street      string    `db:"street" json:"street"`
	City        string    `db:"city" json:"city"`
	State       string    `db:"state" json:"state"`
	Zip         string    `db:"zip" json:"zip"`
	PhoneNumber string    `db:"phone" json:"phone"`
	Created     time.Time `db:"created" json:"created"`
	Modified    time.Time `db:"modified" json:"modified"`
}

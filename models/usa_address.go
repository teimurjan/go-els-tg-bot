package models

import "time"

// UsaAddress is a model of ELS usa address
type UsaAddress struct {
	ID          int64     `db:"id" json:"id" diff:"-"`
	Street      string    `db:"street" json:"street" diff:"street"`
	City        string    `db:"city" json:"city" diff:"city"`
	State       string    `db:"state" json:"state" diff:"state"`
	Zip         string    `db:"zip" json:"zip" diff:"zip"`
	PhoneNumber string    `db:"phone" json:"phone" diff:"phone"`
	Created     time.Time `db:"created" json:"created" diff:"-"`
	Modified    time.Time `db:"modified" json:"modified" diff:"-"`
}

package models

import "time"

type Tracking struct {
	ID       int64     `db:"id" json:"id"`
	UserID   int64     `db:"user_id" json:"userId"`
	Name     string    `db:"name" json:"name"`
	Value    string    `db:"value" json:"value"`
	Status   string    `db:"status" json:"status"`
	Created  time.Time `db:"created" json:"created"`
	Modified time.Time `db:"modified" json:"modified"`
}

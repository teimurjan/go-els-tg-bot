package models

import "time"

type User struct {
	ID       int64     `db:"id" json:"id"`
	ChatID   int64     `db:"chat_id" json:"chatId"`
	Created  time.Time `db:"created" json:"created"`
	Modified time.Time `db:"modified" json:"modified"`
}

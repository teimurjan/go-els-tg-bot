package models

import "time"

// AddTrackingDialog is a model for monitoring the process of a tracking addition
type AddTrackingDialog struct {
	ID                 int64     `db:"id" json:"id"`
	UserID             int64     `db:"user_id" json:"userId"`
	Step               int64     `db:"step" json:"step"`
	FutureTrackingName string    `db:"future_tracking_name" json:"futureTrackingName"`
	Created            time.Time `db:"created" json:"created"`
	Modified           time.Time `db:"modified" json:"modified"`
}

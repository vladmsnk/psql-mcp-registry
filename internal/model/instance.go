package model

import (
	"time"
)

type Instance struct {
	ID              int       `db:"id"`
	Name            string    `db:"name"`
	DatabaseName    string    `db:"database_name"`
	Description     string    `db:"description"`
	CreatorUsername string    `db:"creator_username"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

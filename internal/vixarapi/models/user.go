package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// User represents a user in database
type User struct {
	ID        string
	Username  string
	Email     string
	TgUser    pgtype.Text
	VKID      pgtype.Int8
	CreatedAt time.Time
}

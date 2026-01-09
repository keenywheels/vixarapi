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

// UserQuery represents a user search query in database
type UserQuery struct {
	ID        string
	UserID    string
	Query     string
	CreatedAt time.Time
}

// UserTokenSub represents a user token subscription in database
type UserTokenSub struct {
	ID               string
	UserID           string
	Token            string
	Category         string
	CurrentInterest  float64
	PreviousInterest float64
	Threshold        float64
	Method           string
	ScanDate         time.Time
	CreatedAt        time.Time
}

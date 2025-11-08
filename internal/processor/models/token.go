package models

import "time"

// TokenData represents the token data in postgres database
type TokenData struct {
	TokenID   int64
	TokenName string
	Interest  int64
	Sentiment int16
	SiteName  string
	Date      time.Time
}

package models

import "time"

// TokenData represents the token data in postgres database
type TokenData struct {
	TokenID   int64
	TokenName string
	Interest  int
	Context   string
	SiteName  string
	Date      time.Time
}

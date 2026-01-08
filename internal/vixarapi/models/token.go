package models

import "time"

// Token represent unparsed token record in database
type Token struct {
	TokenName      string
	ScrapeDate     time.Time
	Interest       int64
	Sentiment      int16
	Category       string
	GlobalMedian   int64
	CategoryMedian int64
}

// TokenRecord represent a single record of token data
type TokenRecord struct {
	ScrapeDate       time.Time
	Interest         int64
	GlobalInterest   float64
	CategoryInterest float64
	Sentiment        int16
}

// TokenInfo represent information about a token in database
type TokenInfo struct {
	TokenName string
	Category  string
	Records   []TokenRecord
}

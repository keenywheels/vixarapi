package repository

import "github.com/keenywheels/backend/pkg/postgres"

// TokenDataFields represents the fields of the token data table
type TokenDataFields struct {
	TokenID   string
	TokenName string
	Interest  string
	Sentiment string
	SiteName  string
	Date      string
}

// TokenDataTable represents the structure of the token data table
type TokenDataTable struct {
	Name   string
	Fields TokenDataFields
}

// Repository struct for repository layer
type Repository struct {
	tbl TokenDataTable
	db  *postgres.Postgres
}

// New creates a new Repository instance
func New(db *postgres.Postgres) *Repository {
	tbl := TokenDataTable{
		Name: "token_data",
		Fields: TokenDataFields{
			TokenID:   "token_id",
			TokenName: "token_name",
			Interest:  "interest",
			Sentiment: "sentiment",
			SiteName:  "site_name",
			Date:      "scrape_date",
		},
	}

	return &Repository{
		tbl: tbl,
		db:  db,
	}
}

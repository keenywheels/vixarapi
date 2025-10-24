package models

// Features represent features of an interest instance
type Features struct {
	Interest int `json:"interest"`
}

// Interest represent an interest record
type Interest struct {
	Timestamp int64
	Features  Features // saved as JSONB in database
}

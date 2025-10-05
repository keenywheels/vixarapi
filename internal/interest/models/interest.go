package models

// Features represent features of an interest instance
type Features struct {
	Interest int
}

// Interest represent an interest record
type Interest struct {
	Timestamp int64
	Features  Features
}

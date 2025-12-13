package postgres

// SearchTokenFields represents the fields of the search token table
type SearchTokenFields struct {
	TokenName      string
	ScrapeDate     string
	Interest       string
	Sentiment      string
	MedianInterest string
}

// SearchTokenTable represents the structure of the search token table
type SearchTokenTable struct {
	Name   string
	Fields SearchTokenFields
}

// UserFields represents the fields of the user table
type UserFields struct {
	ID        string
	Username  string
	Email     string
	VKID      string
	TgUser    string
	CreatedAt string
}

// UserTable represents the structure of the user table
type UserTable struct {
	Name   string
	Fields UserFields
}

// UserQueryTable represents the structure of the user query table
type UserQueryTable struct {
	Name   string
	Fields UserQueryFields
}

// UserQueryFields represents the fields of the user query table
type UserQueryFields struct {
	ID        string
	UserID    string
	Query     string
	CreatedAt string
}

// Tables holds all table definitions
type Tables struct {
	search    SearchTokenTable
	user      UserTable
	userQuery UserQueryTable
}

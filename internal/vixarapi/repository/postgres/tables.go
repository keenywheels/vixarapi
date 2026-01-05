package postgres

// SearchTokenFields represents the fields of the search token table
type SearchTokenFields struct {
	TokenName      string
	ScrapeDate     string
	Interest       string
	Sentiment      string
	Category       string
	GlobalMedian   string
	CategoryMedian string
}

// SearchTokenTable represents the structure of the search token table
type SearchTokenTable struct {
	Name   string
	Fields SearchTokenFields
}

// NewSearchTokenTable creates a new instance of SearchTokenTable
func NewSearchTokenTable() SearchTokenTable {
	return SearchTokenTable{
		Name: "mv_token_search",
		Fields: SearchTokenFields{
			TokenName:      "token_name",
			ScrapeDate:     "scrape_date",
			Interest:       "interest",
			Sentiment:      "sentiment",
			Category:       "category",
			GlobalMedian:   "global_median",
			CategoryMedian: "category_median",
		},
	}
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

// NewUserTable creates a new instance of UserTable
func NewUserTable() UserTable {
	return UserTable{
		Name: "users",
		Fields: UserFields{
			ID:        "id",
			Username:  "username",
			Email:     "email",
			TgUser:    "tguser",
			VKID:      "vkid",
			CreatedAt: "created_at",
		},
	}
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

// NewUserQueryTable creates a new instance of UserQueryTable
func NewUserQueryTable() UserQueryTable {
	return UserQueryTable{
		Name: "user_query",
		Fields: UserQueryFields{
			ID:        "id",
			UserID:    "user_id",
			Query:     "query",
			CreatedAt: "created_at",
		},
	}
}

package vk

// ExchangeCodeToTokensResponse response from ExchangeCodeToTokens
type ExchangeCodeToTokensResponse struct {
	ErrorResponse
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	UserID       int64  `json:"user_id"`
	State        string `json:"state"`
	Scope        string `json:"scope"`
}

// RefreshTokensResponse response from RefreshTokens
type RefreshTokensResponse struct {
	ErrorResponse
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	UserID       int64  `json:"user_id"`
	State        string `json:"state"`
	Scope        string `json:"scope"`
}

// UserInfoResponse response from GetUserInfo
type UserInfoResponse struct {
	ErrorResponse
	User struct {
		UserID    string `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Avatar    string `json:"avatar"`
		Email     string `json:"email"`
		Sex       int    `json:"sex"`
		Verified  bool   `json:"verified"`
		Birthday  string `json:"birthday"`
	} `json:"user"`
}

// ErrorResponse VK API error response
type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	State            string `json:"state"`
}

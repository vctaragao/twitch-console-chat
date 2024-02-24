package entity

type AuthToken struct {
	ID           int      `json:"id"`
	Username     string   `json:"username"`
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"`
	Scopes       []string `json:"scopes"`
}

package jwt

type Command struct {
	Id           string
	UserID       string
	AccessToken  string
	RefreshToken string
}

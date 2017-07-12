package entities

type User struct {
	Id string `json,gorethink:"id,omitempty"`
	Username string `json,gorethink:"username"`
	Password string
	Role Role `json,gorethink:"role"`
	Token Token
}

type Token struct {
	Token string `json,gorethink:"token"`
	Expire int64 `json,gorethink:"expire"`
	RefreshToken string `json,gorethink:"refresh_token"`
}

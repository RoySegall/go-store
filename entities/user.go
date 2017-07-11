package entities

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

type User struct {
	Id string `json:"id,omitempty"`
	Username string `json:"title"`
	Password string
	Role Role `json:"role"`
	Token string
}

type Token struct {
	Token string `json:"token"`
	Expire int64 `json:"expire"`
	RefreshToken string `json:"refresh_token"`
}

// Generating a token.
func generateToken(mode string, user User) (string) {
	token_template := jwt.MapClaims {
		"name": user.Username,
		"mode": mode,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, token_template)
	signKey, err := ioutil.ReadFile("private.key")

	if err != nil {
		panic(err)
	}

	tokenResult, err := t.SignedString(signKey)

	if err != nil {
		panic(err)
	}

	return tokenResult
}

// Generating token.
func (token *Token) Generate(user User) () {
	token.Token = generateToken("main_token", user)
	token.Expire = time.Now().Add(time.Hour * 24).Unix()
	token.RefreshToken = generateToken("refresh_token", user)
}

// Checking if the token is valid or not.
func (token Token) Validate() (bool) {

	if token.Expire < time.Now().Unix() {
		// This is an expired token.
		return false
	}

	parsingResult, err := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
		signKey, err := ioutil.ReadFile("private.key")

		if err != nil {
			panic(err)
		}

		return []byte(signKey), nil
	})

	if err != nil {
		return false
	}

	return parsingResult.Valid
}

package entities

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"fmt"
)

type User struct {
	Id string `json:"id,omitempty"`
	Username string `json:"title"`
	Password string `json:"price"`
	Role Role `json:"role"`
	Token Token
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
		"pass": user.Password,
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

	validation_token, err := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return ok, nil
	})

	return true
}

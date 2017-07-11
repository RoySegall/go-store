package entities

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"store/api"
	"go/token"
)

type User struct {
	Id string `gorethink:"id,omitempty"`
	Username string `json:"title"`
	Password string
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

func (user User) Insert() (User, error) {

	// Check if the user exists in the DB.
	if true {
		return nil, error("a")
	}


	// Generating the token object.
	token := Token{}
	token.Generate(user)
	user.Token = token

	// Set the role.
	user.Role = Role{
		Title: "Member",
	}

	// Insert into the DB.
	id := api.Insert("user", user)
	user.Id = id

	return user, nil
}

// Register end point.
func UserRegister(w http.ResponseWriter, r *http.Request) {
	// Get a user input.
	user := User{}
	json.NewDecoder(r.Body).Decode(&user)

	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	// Construct a safe object.
	clean_user := User {
		Username: user.Username,
		Password: string(bytes),
	}

	// Create the user.
	clean_user, err := clean_user.Insert()

	// Print the items.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response []byte
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ = json.Marshal(map[string] error {
			"data": err,
		})
	} else {
		response, _ = json.Marshal(map[string] User {
			"data": clean_user,
		})
	}

	w.Write(response)
}
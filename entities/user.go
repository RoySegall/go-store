package entities

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	r "gopkg.in/gorethink/gorethink.v3"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"store/api"
	"errors"
	"log"
	"fmt"
)

type User struct {
	Id string `gorethink:"id,omitempty"`
	Username string `json:"username"`
	Password string
	Role Role `json:"role"`
	Token Token
}

type Token struct {
	Token string `json:"token"`
	Expire int64 `json:"expire"`
	RefreshToken string `json:"refresh_token"`
}

// Display an error easily.
func WriteError(writer http.ResponseWriter, errorMessage string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)

	response, _ := json.Marshal(map[string] string {
		"data": errorMessage,
	})

	writer.Write(response)
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
	res, err := r.DB("store").Table("user").Filter(map[string]interface{} {
		"Username": user.Username,
	}).Run(api.GetSession())

	if err != nil {
		log.Print(err)
		return user, errors.New("There was an error. Please try again.")
	}

	users := []User{}
	res.All(&users)

	if len(users) != 0 {
		return user, errors.New("User already exists. Try another one.")
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
func UserRegister(writer http.ResponseWriter, request *http.Request) {
	// Get a user input.
	user := User{}
	json.NewDecoder(request.Body).Decode(&user)

	if user.Username == "" || len(user.Username) < 3 {
		WriteError(writer, "Username cannot be empty or less than 3 characters")
		return
	}

	if user.Password == "" || len(user.Password) < 6 {
		WriteError(writer, "Password cannot be empty or less than 6 characters")
		return
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	// Construct a safe object.
	clean_user := User {
		Username: user.Username,
		Password: string(bytes),
	}

	// Create the user.
	clean_user, err := clean_user.Insert()

	// Print the items.
	if err != nil {
		s := err.Error()
		WriteError(writer, s)
		return
	}

	response, _ := json.Marshal(map[string] User {
		"data": clean_user,
	})

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

// Login the user.
func UserLogin(writer http.ResponseWriter, request *http.Request) {
	// Get a user input.
	user := User{}
	json.NewDecoder(request.Body).Decode(&user)

	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	// Construct a safe object.
	clean_user := User {
		Username: user.Username,
		Password: string(bytes),
	}

	res, err := r.DB("store").Table("user").Filter(map[string]interface{} {
		"Username": clean_user.Username,
		//"Password": clean_user.Password,
	}).Run(api.GetSession())

	// Check if the username exists.
	if err != nil {
		log.Print(err)
		WriteError(writer, "The password and the user are wrong. Try again please.")
		return
	}

	users := []User{}
	res.All(&users)
	fmt.Println(res)
	fmt.Println(clean_user)

	//users := []User{}
	//res.All(&users)
	//
	//response, _ := json.Marshal(map[string] User {
	//	"data": users[0],
	//})
	//
	//writer.Header().Set("Content-Type", "application/json")
	//writer.WriteHeader(http.StatusOK)
	//writer.Write(response)
}

// Refreshing an old token.
func UserTokenRefresh(writer http.ResponseWriter, request *http.Request) {

}

// Update user details.
func UserUpdate(writer http.ResponseWriter, request *http.Request) {

}

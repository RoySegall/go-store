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
)

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

// Generating a token.
func generateToken(mode string, user User) (string) {
	token_template := jwt.MapClaims {
		"name": user.Username,
		"mode": mode,
		"time": int64(time.Now().Unix()) + int64(time.Now().Nanosecond()),
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

// Load a user from the DB with a token.
func LoadUserFromDB(token string) (User, error) {

	user := User{}

	// Query from the DB.
	res, err := r.DB("store").Table("user").Filter(map[string]interface{} {
		"Token": map[string]interface{} {
			"Token": token,
		},
	}).Run(api.GetSession())

	if err != nil {
		s := err.Error()
		log.Print(s)
		return user, errors.New("There was an error.")
	}

	users := []User{}
	res.All(&users)

	if len(users) == 0 {
		return user, errors.New("There is no user a matching access token.")
	}

	if !users[0].Token.Validate() {
		return user, errors.New("The token is not valid. Might be expired.")
	}

	return users[0], nil
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

// Update the user object.
func (user User) Update() {
	api.Update("user", user)
}

// Register end point.
func UserRegister(writer http.ResponseWriter, request *http.Request) {
	// Get a user input.
	user := User{}
	json.NewDecoder(request.Body).Decode(&user)

	if user.Username == "" || len(user.Username) < 3 {
		api.WriteError(writer, "Username cannot be empty or less than 3 characters")
		return
	}

	if user.Password == "" || len(user.Password) < 6 {
		api.WriteError(writer, "Password cannot be empty or less than 6 characters")
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
		api.WriteError(writer, s)
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
	user := &User{}
	json.NewDecoder(request.Body).Decode(&user)

	res, err := r.DB("store").Table("user").Filter(map[string]interface{} {
		"Username": user.Username,
	}).Run(api.GetSession())

	// Check if the username exists.
	if err != nil {
		log.Print(err)
		api.WriteError(writer, "The password and the user are wrong. Try again please.")
		return
	}

	// Prepare the result from the DB.
	DbUsers := []User{}
	res.All(&DbUsers)
	matchedUser := DbUsers[0]

	err = bcrypt.CompareHashAndPassword([]byte(matchedUser.Password), []byte(user.Password))

	if err != nil {
		api.WriteError(writer, "The password and the user are wrong. Try again please.")
		return
	}

	// Create a new token for the user.
	token := Token{}
	token.Generate(matchedUser)
	matchedUser.Token = token
	matchedUser.Update()

	// Display the user.
	response, _ := json.Marshal(map[string] User {
		"data": matchedUser,
	})

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

// Refreshing an old token.
func UserTokenRefresh(writer http.ResponseWriter, request *http.Request) {
	// Get a user input.
	token := &Token{}
	json.NewDecoder(request.Body).Decode(&token)

	if token.RefreshToken == "" {
		// A token cannot be empty. Notify the client about that.
		api.WriteError(writer, "The refresh token cannot be empty")
		return
	}

	// Query from the DB.
	res, err := r.DB("store").Table("user").Filter(map[string]interface{} {
		"Token": map[string]interface{} {
			"RefreshToken": token.RefreshToken,
		},
	}).Run(api.GetSession())

	if err != nil {
		log.Print(err)
		api.WriteError(writer, "It seems there was an error. Try again later")
		return
	}

	// Prepare the result from the DB.
	DbUsers := []User{}
	res.All(&DbUsers)

	if len(DbUsers) == 0 {
		api.WriteError(writer, "There is no user with that refresh token. Try again.")
		return
	}

	matchedUser := DbUsers[0]

	// Create a new token for the user.
	new_token := Token{}
	new_token.Generate(matchedUser)
	matchedUser.Token = new_token
	matchedUser.Update()

	// Display the user.
	response, _ := json.Marshal(map[string] User {
		"data": matchedUser,
	})

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

// Get user details.
func UserInfo(writer http.ResponseWriter, request *http.Request) {
	token := request.Header.Get("access-token")
	user, err := LoadUserFromDB(token)

	if err != nil {
		s := err.Error()
		api.WriteError(writer, s)
		return
	}

	// Display the user.
	response, _ := json.Marshal(map[string] User {
		"data": user,
	})

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}


package entities

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	r "gopkg.in/gorethink/gorethink.v3"
	"io/ioutil"
	"store/api"
	"errors"
	"log"
)

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

// Insert a user into the DB.
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
	api.Update("user", user.Id, user)
}

// Load user from DB.
func (user User) Get(id string) (User, error) {
	// Check if the user exists in the DB.
	res, err := r.DB("store").Table("user").Filter(map[string]interface{} {
		"id": id,
	}).Run(api.GetSession())

	if err != nil {
		log.Print(err)
		return user, errors.New("There was an error. Please try again.")
	}

	users := []User{}
	res.All(&users)

	if len(users) == 0 {
		return user, errors.New("User does not exists. Try another one.")
	}

	return users[0], nil
}

// Adding an item to the user's cart.
func (user *User) AddItemToCart(item Item) {
	// Append the item to the cart property.
	user.Cart.Items = append(user.Cart.Items, item)
}

// Revoke an item from the cart.
func (user *User) RevokeItemFromCart(itemId string) {

	var items []Item

	for _, item := range user.Cart.Items {

		if item.Id == itemId {
			continue
		}

		items = append(items, item)
	}

	user.Cart.Items = items
}

// Tak the current cart and archive it.
func (user *User) ArchiveCart() {
	user.PastCarts = append(user.PastCarts, user.Cart)
	user.Cart = Cart{}
}

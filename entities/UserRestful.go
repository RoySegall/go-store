package entities

import (
	r "gopkg.in/gorethink/gorethink.v3"
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"store/api"
	"log"
)

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

// Adding an item to the user cart.
func UserAddItemToCart(write http.ResponseWriter, request *http.Request) {

}

// Removing an item from the cart.
func UserRevokeItemFromCart(write http.ResponseWriter, request *http.Request) {

}

// After the user finished with the cart, archive it.
func UserArchiveCart(write http.ResponseWriter, request *http.Request) {

}

package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"store/entities"
	"github.com/fatih/color"
)

func main() {

	r := mux.NewRouter()

	// Items.
	r.HandleFunc("/api/item", entities.ItemsGet).Methods(http.MethodGet)
	r.HandleFunc("/api/item", entities.ItemPost).Methods(http.MethodPost)
	r.HandleFunc("/api/item/{id}", entities.ItemGet).Methods(http.MethodGet)
	r.HandleFunc("/api/item/{id}", entities.ItemUpdate).Methods(http.MethodPatch)
	r.HandleFunc("/api/item/{id}", entities.ItemDelete).Methods(http.MethodDelete)

	// User.
	r.HandleFunc("/api/user", entities.UserInfo).Methods(http.MethodGet)
	r.HandleFunc("/api/user", entities.UserRegister).Methods(http.MethodPost)
	r.HandleFunc("/api/user/login", entities.UserLogin).Methods(http.MethodPost)
	r.HandleFunc("/api/user/token_refresh", entities.UserTokenRefresh).Methods(http.MethodPost)

	// Cart management.
	r.HandleFunc("/api/cart/items", entities.UserAddItemToCart).Methods(http.MethodPost)
	r.HandleFunc("/api/cart/items", entities.UserRevokeItemFromCart).Methods(http.MethodDelete)
	r.HandleFunc("/api/cart", entities.UserArchiveCart).Methods(http.MethodDelete)

	color.Green("Server started...")

	server := &http.Server{
		Addr: ":8070",
		Handler: r,
	}

	server.ListenAndServe()
}

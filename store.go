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
	r.HandleFunc("/api/item", entities.ItemsGet).Methods("GET")
	r.HandleFunc("/api/item", entities.ItemPost).Methods("POST")
	r.HandleFunc("/api/item/{id}", entities.ItemGet).Methods("GET")
	r.HandleFunc("/api/item/{id}", entities.ItemUpdate).Methods("PATCH")
	r.HandleFunc("/api/item/{id}", entities.ItemDelete).Methods("DELETE")

	// User.
	r.HandleFunc("/api/user", entities.UserRegister).Methods("POST")

	color.Green("Server started...")

	server := &http.Server{
		Addr: ":8070",
		Handler: r,
	}

	server.ListenAndServe()
}

package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"store/entities"
	"github.com/fatih/color"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/item", entities.ItemsGet).Methods("GET")
	r.HandleFunc("/item", entities.ItemPost).Methods("POST")
	r.HandleFunc("/item/{id}", entities.ItemGet).Methods("GET")
	r.HandleFunc("/item/{id}", entities.ItemUpdate).Methods("PUT")
	r.HandleFunc("/item/{id}", entities.ItemDelete).Methods("DELETE")

	color.Green("Server started...")
	http.ListenAndServe(":8070", r)
}

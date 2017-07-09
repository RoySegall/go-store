package entities

import (
	"store/api"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

type Item struct {
	Id string
	Title string
	Price float64
}

// Insert an object.
func (item Item) Insert() (string) {
	return api.Insert("item", item)
}

// Get a single object.
func (item Item) Get(id string) (Item) {
	res := api.Get("item", id)
	items := []Item{}
	res.All(&items)
	return items[0]
}

// Get all the items.
func (item Item) GetAll() ([]Item) {
	res := api.GetAll("item")
	items := []Item{}
	res.All(&items)
	return items
}

// Delete an item.
func (item Item) Delete() {
	api.Delete("item", item.Id)
}

func ItemsGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Hello!\n")
}

func ItemGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["id"])
}

func ItemPost(w http.ResponseWriter, r *http.Request) {

}

func ItemUpdate(w http.ResponseWriter, r *http.Request) {

}

func ItemDelete(w http.ResponseWriter, r *http.Request) {

}

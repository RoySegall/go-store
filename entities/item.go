package entities

import (
	"store/api"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
)

type Item struct {
	Id 		string 	`json:"id,omitempty"`
	Title string 	`json:"title"`
	Price float64 `json:"price"`
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

// Get a specific item.
func ItemsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response, _ := json.Marshal(Item{}.GetAll())

	w.Write(response)
}

// Get a specific item.
func ItemGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	response, _ := json.Marshal(Item{}.Get(vars["id"]))
	w.Write(response)
}

func ItemPost(w http.ResponseWriter, r *http.Request) {

}

func ItemUpdate(w http.ResponseWriter, r *http.Request) {

}

func ItemDelete(w http.ResponseWriter, r *http.Request) {

}

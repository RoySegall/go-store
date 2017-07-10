package entities

import (
	"store/api"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
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

	// Print out the
	response, _ := json.Marshal(map[string] []Item {
		"data": Item{}.GetAll(),
	})

	// Print the items.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Get a specific item.
func ItemGet(w http.ResponseWriter, r *http.Request) {
	// Pull a single item from the DB.
	vars := mux.Vars(r)
	response, _ := json.Marshal(map[string] Item {
		"data": Item{}.Get(vars["id"]),
	})

	// Print the items.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func ItemPost(w http.ResponseWriter, r *http.Request) {
	// Processing.
	item := Item{}
	json.NewDecoder(r.Body).Decode(&item)
	id := item.Insert()

	// Adding the ID to the object.
	item.Id = id

	// Prepare the display.
	response, _ := json.Marshal(map[string] Item {
		"data": item,
	})

	// Print, with style.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func ItemUpdate(w http.ResponseWriter, r *http.Request) {

}

func ItemDelete(w http.ResponseWriter, r *http.Request) {

}

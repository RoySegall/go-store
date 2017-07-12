package entities

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"store/api"
	"log"
	"io/ioutil"
	"github.com/imdario/mergo"
	"encoding/base64"
)

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

	if item.Title == "" {
		api.WriteError(w, "The title is required")
		return
	}

	if item.Price == 0 {
		api.WriteError(w, "An item price cannot be 0")
		return
	}

	// Before creating the entry in the DB, we need to save the image.
	if item.Image == "" {
		api.WriteError(w, "You need to provide an image")
		return
	}

	buff, err := base64.StdEncoding.DecodeString(item.Image)

	if err != nil {
		s := err.Error()
		log.Print(err)
		api.WriteError(w, s)
		return
	}

	settings := api.GetSettings()
	if err := ioutil.WriteFile(settings.ImageDirectory + "/" + item.Title + ".jpg", buff, 777); err != nil {
		s := err.Error()
		log.Print(err)
		api.WriteError(w, s)
		return
	}

	item.Image = item.Title + ".jpg"
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

// Update an item.
func ItemUpdate(w http.ResponseWriter, r *http.Request) {
	// Process variables.
	vars := mux.Vars(r)

	// Old object.
	old_item := Item{}.Get(vars["id"])

	// Process the new values and attach the ID to the object.
	item := Item{}
	json.NewDecoder(r.Body).Decode(&item)

	if item.Title == "" {
		api.WriteError(w, "The title is required")
		return
	}

	if item.Price == 0 {
		api.WriteError(w, "An item price cannot be 0")
		return
	}

	if err := mergo.Merge(&item, old_item); err != nil {
		log.Print(err)
		api.WriteError(w, "It seems that was an error. Try again later")
		return
	}

	// Updating.
	item.Update()

	// Prepare the display.
	response, _ := json.Marshal(map[string] Item {
		"data": item,
	})

	// Print, with style.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Delete an item.
func ItemDelete(w http.ResponseWriter, r *http.Request) {
	// Process variables.
	vars := mux.Vars(r)

	// Delete the item.
	Item{}.Get(vars["id"]).Delete()

	response, _ := json.Marshal(map[string] string {
		"result": "deleted",
	})

	// Print, with style.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

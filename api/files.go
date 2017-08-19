package api

import (
	"net/http"
	"os"
	"log"
	"io"
	"github.com/gorilla/mux"
)

// Serving a file.
func ServeFile(w http.ResponseWriter, r *http.Request) {
	// Pull a single item from the DB.
	vars := mux.Vars(r)

	img, err := os.Open("./images/" + vars["file"])
	if err != nil {
		log.Fatal(err) // perhaps handle this nicer
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, img)
}

package entities

import (
	"net/http"
	"github.com/labstack/echo"
	"os"
	"io"
	"store/api"
)

// Get a specific item.
func ItemsGet(c echo.Context) error {

	// Print out the
	return c.JSON(200, map[string] []Item {
		"data": Item{}.GetAll(),
	})
}

// Get a specific item.
func ItemGet(c echo.Context) error {
	// Print the items.
	return c.JSON(200,	map[string] Item {
		"data": Item{}.Get(c.Param("id")),
	})
}

// Posting an item.
func ItemPost(c echo.Context) error {
	// Processing.
	item := new(Item)

	if err := c.Bind(item); err != nil {
		return err
	}

	if item.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "The title is required.")
	}

	if item.Price == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "An item price cannot be 0")
	}

	// Before creating the entry in the DB, we need to save the image.
	file, err := c.FormFile("image")
	if file == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "You need to provide an image")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	settings := api.GetSettings()
	dst, err := os.Create(settings.ImageDirectory + file.Filename)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	defer dst.Close()

	// Copy.
	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	item.Image = settings.ImageDirectory + file.Filename
	id := item.Insert()

	// Adding the ID to the object.
	item.Id = id

	// Prepare the display.
	return c.JSON(200,	item)
}

//// Update an item.
//func ItemUpdate(c echo.Context) error {
//	// Process variables.
//	vars := mux.Vars(r)
//
//	// Old object.
//	old_item := Item{}.Get(vars["id"])
//
//	// Process the new values and attach the ID to the object.
//	item := Item{}
//	json.NewDecoder(r.Body).Decode(&item)
//
//	if item.Title == "" {
//		api.WriteError(w, "The title is required")
//		return
//	}
//
//	if item.Price == 0 {
//		api.WriteError(w, "An item price cannot be 0")
//		return
//	}
//
//	if err := mergo.Merge(&item, old_item); err != nil {
//		log.Print(err)
//		api.WriteError(w, "It seems that was an error. Try again later")
//		return
//	}
//
//	// Updating.
//	item.Update()
//
//	// Prepare the display.
//	response, _ := json.Marshal(map[string] Item {
//		"data": item,
//	})
//
//	// Print, with style.
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	w.Write(response)
//}
//
//// Delete an item.
//func ItemDelete(c echo.Context) error {
//	// Process variables.
//	vars := mux.Vars(r)
//
//	// Delete the item.
//	Item{}.Get(vars["id"]).Delete()
//
//	response, _ := json.Marshal(map[string] string {
//		"result": "deleted",
//	})
//
//	// Print, with style.
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	w.Write(response)
//}

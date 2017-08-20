package entities

import (
	"net/http"
	"github.com/labstack/echo"
	"os"
	"io"
	"store/api"
	"github.com/imdario/mergo"
	"github.com/labstack/gommon/log"
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

	object := Item{}.Get(c.Param("id"))

	if object.Id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "The item does not exists any more.")
	}

	// Print the items.
	return c.JSON(200,	map[string] Item {
		"data": object,
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

	if item.Description == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "The description is mandatory.")
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

// Update an item.
func ItemUpdate(c echo.Context) error {
	// Old object.
	old_item := Item{}.Get(c.Param("id"))

	// Processing.
	item := new(Item)

	if err := c.Bind(item); err != nil {
		return err
	}

	// todo move to item.Validate()
	if item.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "The title is required.")
	}

	if item.Price == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "An item price cannot be 0.")
	}

	if item.Description == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "The description is mandatory.")
	}

	if err := mergo.Merge(item, old_item); err != nil {
		log.Print(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Updating.
	item.Update()

	// Prepare the display.
	return c.JSON(200,	item)
}

// Delete an item.
func ItemDelete(c echo.Context) error {
	// Delete the item.
	Item{}.Get(c.Param("id")).Delete()

	return c.JSON(http.StatusOK, map[string] string {
		"result": "deleted",
	})
}

package entities

import (
	"github.com/labstack/echo"
	"net/http"
)

// Adding an item to the user cart.
func UserAddItemToCart(c echo.Context) error {

	user, err := LoadUserFromDB(c.Request().Header.Get("access-token"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "There is not matching user. Maybe a bas access token?")
	}

	// Get a user input.
	post_item := c.FormValue("item_id")

	if post_item == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Item cannot be empty!")
	}

	item := Item{}.Get(post_item)

	if item.Id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "There is no item with the ID " + post_item)
	}

	user.AddItemToCart(item)
	user.Update()
	return c.JSON(http.StatusOK, map[string]string{"message": "The item " + item.Title + " Was added to the cart."})
}

// Removing an item from the cart.
func UserRevokeItemFromCart(c echo.Context) error {
	user, err := LoadUserFromDB(c.Request().Header.Get("access-token"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Get a user input.
	post_item := c.FormValue("item_id")

	if post_item == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Item cannot be empty!")
	}

	item := Item{}.Get(post_item)

	if item.Id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "There is no item with the ID " + post_item)
	}

	user.RevokeItemFromCart(post_item)
	user.Update()
	return c.JSON(http.StatusOK, map[string]string{
		"message": "The item " + item.Title + " revoked from the the cart.",
	})
}

// After the user finished with the cart, archive it.
func UserArchiveCart(c echo.Context) error {
	user, err := LoadUserFromDB(c.Request().Header.Get("access-token"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if len(user.Cart.Items) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,"There are no items in the current cart.")
	}

	user.ArchiveCart()
	user.Update()
	return c.JSON(http.StatusOK, map[string]string{
		"message": "The cart moved to the archive",
	})
}

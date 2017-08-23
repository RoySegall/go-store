package api

import (
	"github.com/labstack/echo"
)

// Serving a file.
func ServeFile(c echo.Context) error {
	return c.File("./images/" + c.Param("file"))
}

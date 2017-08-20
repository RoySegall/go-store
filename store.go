package main

import (
	"store/entities"
	"github.com/fatih/color"
	"store/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// Items.
	e.GET("/api/items", entities.ItemsGet)
	e.POST("/api/items", entities.ItemPost)
	e.GET("/api/items/:id", entities.ItemGet)
	e.PATCH("/api/items/:id", entities.ItemUpdate)
	e.DELETE("/api/items/:id", entities.ItemDelete)
	//
	//// User.
	//e.GET("/api/user", entities.UserInfo)
	//e.POST("/api/user", entities.UserRegister)
	e.POST("/api/user/login", entities.UserLogin)
	//e.POST("/api/user/token_refresh", entities.UserTokenRefresh)
	//
	//// Cart management.
	//e.POST("/api/cart/items", entities.UserAddItemToCart)
	//e.DELETE("/api/cart/items", entities.UserRevokeItemFromCart)
	//e.DELETE("/api/cart", entities.UserArchiveCart)
	//
	//// Handle files.
	//e.GET("/images/{file}", api.ServeFile)

	color.Green("Starting server at http://localhost" + api.GetSettings().Port)

	e.Start(api.GetSettings().Port)
}

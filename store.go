package main

import (
	"store/entities"
	"github.com/fatih/color"
	"store/api"
	"github.com/labstack/echo"
)

func main() {


	//
	//// User.
	//r.HandleFunc("/api/user", entities.UserInfo).Methods(http.MethodGet)
	//r.HandleFunc("/api/user", entities.UserRegister).Methods(http.MethodPost)
	//r.HandleFunc("/api/user/login", entities.UserLogin).Methods(http.MethodPost)
	//r.HandleFunc("/api/user/token_refresh", entities.UserTokenRefresh).Methods(http.MethodPost)
	//
	//// Cart management.
	//r.HandleFunc("/api/cart/items", entities.UserAddItemToCart).Methods(http.MethodPost)
	//r.HandleFunc("/api/cart/items", entities.UserRevokeItemFromCart).Methods(http.MethodDelete)
	//r.HandleFunc("/api/cart", entities.UserArchiveCart).Methods(http.MethodDelete)
	//
	//// Handle files.
	//r.HandleFunc("/images/{file}", api.ServeFile).Methods(http.MethodGet)

	e := echo.New()

	//Items.
	e.GET("/api/items", entities.ItemsGet)
	//e.POST("/api/item", entities.ItemPost)
	//e.GET("/api/item/{id}", entities.ItemGet)
	//e.PATCH("/api/item/{id}", entities.ItemUpdate)
	//e.DELETE("/api/item/{id}", entities.ItemDelete)

	color.Green("Starting server at http://localhost" + api.GetSettings().Port)
	e.Logger.Fatal(e.Start(api.GetSettings().Port))
}

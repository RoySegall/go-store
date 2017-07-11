package main

import (
	"github.com/fatih/color"
	"store/api"
)

func main() {
	color.Yellow("Starting to install.")
	api.DbCreate("store")
	color.Green("The DB 'store' has created.")

	color.Yellow("Starting to create tables.")

	for _, table := range []string{"item", "users"} {
		api.TableCreate(table)
		color.Green("The %s table has created.\n", table)
	}

	color.GreenString("Done!")
}

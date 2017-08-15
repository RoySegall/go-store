package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"store/api"
)

func main() {

	if !askQuestion("Do you want to start? This will erase any content!") {
		fmt.Println("Quitting.")
		return
	}

	if askQuestion("Migrate users?") {
		userMigrate()
	}

	if askQuestion("Migrate items") {
		itemsMigrate()
	}
}

// Asking a question.
func askQuestion(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(question + " Y/n")

	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	if !api.InArray([]string{"y", "Y", "n", "N"}, option) {
		fmt.Println("Not a valid option. Skipping")
		return false
	}

	if api.InArray([]string{"N", "n"}, option) {
		return false
	}

	return true
}

func userMigrate() {
	fmt.Println("migratin users")
}

func itemsMigrate() {
	fmt.Println("migratin items")
}

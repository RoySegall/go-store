package main

import (
	"bufio"
	"os"
	"fmt"
	"strings"
	"store/api"
	"path/filepath"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"store/entities"
	"log"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	userMigrate()
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
	fmt.Println("Drop users table")
	api.TableDrop("user")
	fmt.Println("Create users table")
	api.TableCreate("user")

	fmt.Println("migratin users")

	filename, _ := filepath.Abs("./migration/users.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	m := make(map[string]entities.User)

	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, user := range m {
		// Alter the user password.
		bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		user.Password = string(bytes)

		// Migrate the image.
		//dest_image_path, _ := filepath.Abs(api.GetSettings().ImageDirectory)
		//source_image_path, _ := filepath.Abs("./migration/images/")

		if err != nil {
			panic(err)
		}

		user.Image = api.GetSettings().ImageDirectory + "users/" + user.Image

		object, err := user.Insert()

		if err != nil {
			panic(err)
		}

		fmt.Println("Migrating " + object.Username)
	}
}

func itemsMigrate() {
	fmt.Println("migratin items")
}

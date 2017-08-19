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
	"io"
)

func main() {
	if !askQuestion("Starting the migration process?") {
		fmt.Print("Thnaks god I asked! Your DB good go away!")
		return
	}

	if askQuestion("Migrate users?") {
		userMigrate()
	}

	if askQuestion("Migrate items?") {
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

	fmt.Println("Drop users table")
	api.TableDrop("user")

	fmt.Println("Create users table")
	api.TableCreate("user")

	fmt.Println("migratin users")

	filename, _ := filepath.Abs("./migration_assets/users.yml")
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
		err := fileCopy(
			"migration_assets/images/users/" + user.Image,
			api.GetSettings().ImageDirectory + "/" + user.Image)

		if err != nil {
			panic(err)
		}

		user.Image = strings.Replace(api.GetSettings().ImageDirectory, "./", "", -1) + user.Image

		object, err := user.Insert()

		if err != nil {
			panic(err)
		}

		fmt.Println("Migrating " + object.Username)
	}

	fmt.Println("All the users have been migrated. Awesome!")
}

// Migrating items.
func itemsMigrate() {
	fmt.Println("migratin items")

	fmt.Println("Drop items table")
	api.TableDrop("item")

	fmt.Println("Create items table")
	api.TableCreate("item")

	filename, _ := filepath.Abs("./migration_assets/items.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	m := make(map[string]entities.Item)

	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, item := range m {
		// Migrate the image.
		err := fileCopy(
			"migration_assets/images/items/" + item.Image,
			api.GetSettings().ImageDirectory + "/" + item.Image)

		if err != nil {
			panic(err)
		}

		item.Image = strings.Replace(api.GetSettings().ImageDirectory, "./", "", -1) + item.Image

		item.Insert()

		if err != nil {
			panic(err)
		}

		fmt.Println("Migrating " + item.Title)
	}

	fmt.Println("All the items have been migrated. Awesome!")
}

// Copy file.
func fileCopy(src, dst string) error {

	data, _ := os.Stat(dst)

	if data != nil {
		// Delete a file which already exists.
		os.Remove(dst)
	}

	s, err := os.Open(src)

	if err != nil {
		return err
	}

	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()

	d, err := os.Create(dst)

	if err != nil {
		return err
	}

	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}

	return d.Close()
}

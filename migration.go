package main

import (
	"bufio"
	"os"
	r "gopkg.in/gorethink/gorethink.v3"
	"strings"
	"store/api"
	"path/filepath"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"store/entities"
	"log"
	"golang.org/x/crypto/bcrypt"
	"io"
	"github.com/fatih/color"
)

func main() {
	if !askQuestion("Starting the migration process?") {
		color.Red("Thanks god I asked! Your DB could go away!")
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

	color.Green(question + " Y/n")

	option, _ := reader.ReadString('\n')
	option = strings.TrimSpace(option)

	if !api.InArray([]string{"y", "Y", "n", "N"}, option) {
		color.Yellow("Not a valid option. Skipping")
		return false
	}

	if api.InArray([]string{"N", "n"}, option) {
		return false
	}

	return true
}

func userMigrate() {
	color.Red("Truncate users table")
	r.DB("store").Table("user").Delete().Run(api.GetSession())

	color.Green("Migrating users")

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

		color.Yellow("Migrating " + object.Username)
	}

	color.Green("All the users have been migrated. Awesome!")
}

// Migrating items.
func itemsMigrate() {
	color.Red("Truncate items table")
	r.DB("store").Table("item").Delete().Run(api.GetSession())

	color.Yellow("Migrating items")

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

		color.Yellow("Migrating " + item.Title)
	}

	color.Green("All the items have been migrated. Awesome!")
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

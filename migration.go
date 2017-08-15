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

	for _, value := range m {
		fmt.Println(value)
	}
}

func itemsMigrate() {
	fmt.Println("migratin items")
}

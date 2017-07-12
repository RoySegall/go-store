package api

import (
	"net/http"
	"encoding/json"
	"path/filepath"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Permissions struct {
	Item CRUD `yaml:"item"`
}

type CRUD struct {
	Create 	[]string `yaml:"create"`
	Read 		[]string `yaml:"read"`
	Update 	[]string `yaml:"update"`
	Delete 	[]string `yaml:"delete"`
}

type Settings struct {
	RethinkDB struct {
		Address string `yaml:"address"`
		Database string `yaml:"database"`
	} `yaml:"rethinkdb"`
	ImageDirectory string `yaml:"image_directory"`
}

// Display an error easily.
func WriteError(writer http.ResponseWriter, errorMessage string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)

	response, _ := json.Marshal(map[string] string {
		"data": errorMessage,
	})

	writer.Write(response)
}

// Verify if a value exists in an array.
func InArray(haystack []string, needle string) (bool) {

	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}

// Get the permissions.
func GetPermissions() (Permissions) {
	filename, _ := filepath.Abs("./config/permissions.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Permissions

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}

// Get settings.
func GetSettings() (Settings) {
	filename, _ := filepath.Abs("./config/settings.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Settings

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}

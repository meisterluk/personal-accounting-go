package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/meisterluk/personal-accounting.go/types"

	"gopkg.in/yaml.v2"
)

// Application contains all configuration settings
type Application struct {
	User     string           `yaml:"user"`
	Entities map[int64]string `yaml:"entities,flow"`
	Tags     []string         `yaml:"tags,flow"`
}

// Transaction represents a transaction which is stored
// in a XML file and potentially analyzed
type Transaction struct {
	From     int64                      `xml:"source"`
	To       int64                      `xml:"destination"`
	Tags     []string                   `xml:"tags>tag"`
	Datetime types.TransactionTimestamp `xml:"timestamp"`
	Amount   types.AmountValue          `xml:"amount"`
}

type xmlStructure struct {
	Transactions []Transaction `xml:"transaction"`
}

// ReadConfig read a configuration from a given filepath
// and writes the settings to the given Application instance
func ReadConfig(filepath string, app *Application) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(fmt.Sprintf("Could not read configuration file: %s", filepath))
	}
	yaml.Unmarshal(content, &app)
}

func contains(elements []string, element string) bool {
	for _, v := range elements {
		if element == v {
			return true
		}
	}
	return false
}

func printHeader(app *Application) {
	fmt.Println("personal-accounting.go  v0.1")
	fmt.Println()
	fmt.Printf("1. configuration settings for %s have been read\n", app.User)
	fmt.Println("2. I will ask you to provide transaction information")
	fmt.Println()
}

func retrieveEntity(entityID string, app *Application) (int64, error) {
	entity, err := strconv.ParseInt(strings.TrimSpace(entityID), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Invalid value received for entity: %s - %s", entityID, err.Error()))
	}
	for id := range app.Entities {
		if id == entity {
			return entity, nil
		}
	}
	return 0, fmt.Errorf("Entity not found: %d", entity)
}

func appendToXML(app *Application, filepath string, t *Transaction) {
	xmlFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()
	content, _ := ioutil.ReadAll(xmlFile)

	var x xmlStructure
	xml.Unmarshal(content, &x)

	fmt.Printf("%v\n", x)
}

func main() {
	app := Application{}
	ReadConfig("accounting.yml", &app)

	printHeader(&app)
	t := Transaction{}
	//queryTransaction(&t, &app)

	appendToXML(&app, "dataset.xml", &t)
}

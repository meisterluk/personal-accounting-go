package main

import (
	"fmt"
	"io/ioutil"

	"github.com/meisterluk/personal-accounting-go/cli"
	"github.com/meisterluk/personal-accounting-go/storage"
	"github.com/meisterluk/personal-accounting-go/types"
	flag "github.com/ogier/pflag"
	yaml "gopkg.in/yaml.v2"
)

// ReadConfig read a configuration from a given filepath
// and writes the settings to the given Application instance
func ReadConfig(filepath string, app *types.Application) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(fmt.Sprintf("Could not read configuration file: %s", filepath))
	}
	yaml.Unmarshal(content, &app)
}

func printHeader(app *types.Application) {
	fmt.Println("personal-accounting-go  v0.1")
	fmt.Println()
	fmt.Printf("1. configuration settings for %s have been read\n", app.User)
	fmt.Println("2. I will ask you to provide transaction information")
	fmt.Println()
}

func main() {
	configFile := flag.StringP("config", "c", "accounting.yml", "configuration file")
	flag.Parse()

	app := types.Application{}
	ReadConfig(*configFile, &app)

	var store types.Storage
	if app.StorageType == "xml" {
		store = storage.XMLStorage{app.XMLFilepath}
	}

	printHeader(&app)
	err := cli.MainLoop(&app, &store)
	if err != nil {
		panic(err.Error())
	}
}

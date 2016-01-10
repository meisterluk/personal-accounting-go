package cli

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/meisterluk/personal-accounting-go/types"
)

type xmlStructure struct {
	Transactions []types.Transaction `xml:"transaction"`
}

func appendToXML(app *types.Application, filepath string, t *types.Transaction) {
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

func addTransaction(app *types.Application, s *types.Storage, t *types.Transaction) error {
	appendToXML(app, "dataset.xml", t)
	return nil
}

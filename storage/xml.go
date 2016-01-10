package storage

import "github.com/meisterluk/personal-accounting-go/types"

// XMLStorage stores its data in an XML file
type XMLStorage struct {
	XMLFilepath string
}

// Add an transaction to the existing dataset in the XML storage
func (x XMLStorage) Add(app *types.Application, t *types.Transaction) error {
	return nil
}

// List all transactions stored in the XML file
func (x XMLStorage) List(app *types.Application, out chan types.Transaction) error {
	return nil
}

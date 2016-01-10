package types

import "time"

// Application contains all configuration settings
type Application struct {
	User        string           `yaml:"user"`
	StorageType string           `yaml:"storage"`
	XMLFilepath string           `yaml:"xmlfile"`
	Entities    map[int64]string `yaml:"entities,flow"`
	Tags        []string         `yaml:"tags,flow"`
}

// TransactionTimestamp is an XML marshallable extension to time.Time
type TransactionTimestamp struct {
	App      *Application
	DateTime time.Time
}

// AmountValue represents a monetary value
type AmountValue struct {
	prefix uint64 // large value before decimal point
	suffix uint8  // small value after decimal point
}

// Transaction represents a transaction which is stored
// in a XML file and potentially analyzed
type Transaction struct {
	App      Application
	From     int64                `xml:"source"`
	To       int64                `xml:"destination"`
	Tags     []string             `xml:"tags>tag"`
	Datetime TransactionTimestamp `xml:"timestamp"`
	Amount   AmountValue          `xml:"amount"`
}

// Storage is an interface every storage backend conforms to
type Storage interface {
	Add(*Application, *Transaction) error
	List(*Application, chan Transaction) error
}

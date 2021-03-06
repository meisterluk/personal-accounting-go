package types

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NewTransactionTimestamp is a constructor for TransactionTimestamp objects
// using a string representation
func NewTransactionTimestamp(app *Application, strrepr string) (*TransactionTimestamp, error) {
	parsed, ok := time.Parse(time.RFC3339, strrepr)
	if ok != nil {
		return nil, fmt.Errorf("Could not decode timestamp %s", strrepr)
	}
	t := TransactionTimestamp{app, parsed}
	return &t, nil
}

func (t *TransactionTimestamp) String() string {
	return time.Time(t.DateTime).Format(time.RFC3339)
}

// UnmarshalXML unmarshals the transaction time from an XML serialization
func (t *TransactionTimestamp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	value := struct{ D string }{}
	d.DecodeElement(&value, &start)
	tt, err := NewTransactionTimestamp(t.App, value.D)
	if err != nil {
		return err
	}
	*t = *tt
	return nil
}

// MarshalXML marshals the transaction time for XML serialization
func (t *TransactionTimestamp) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(t.String(), start)
	return nil
}

// NewAmountValue constructs AmountValue objects given a string representation
func NewAmountValue(strrepr string) (*AmountValue, error) {
	// detect delimiter
	delim := "."
	if strings.Contains(strrepr, ",") {
		delim = ","
	}

	// split with delimiter
	parts := strings.Split(strings.TrimSpace(strrepr), delim)
	if len(parts) > 2 {
		return nil, fmt.Errorf("Invalid monetary value: expected only one separator in %s", strrepr)
	}

	high, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}

	// only digits? then values after decimal point
	if len(parts) == 1 {
		obj := AmountValue{prefix: high, suffix: 0}
		return &obj, nil
	}

	// one delimiter? than prefix + suffix
	if len(parts) > 1 && len(parts[1]) > 2 {
		return nil, fmt.Errorf("You cannot use more than two digits after the decimal points")
	}

	low, err := strconv.ParseUint(parts[1], 10, 8)
	if err != nil {
		return nil, err
	}

	obj := AmountValue{prefix: high, suffix: uint8(low)}
	return &obj, nil
}

func (a *AmountValue) String() string {
	return fmt.Sprintf("%d.%d", a.prefix, a.suffix)
}

// UnmarshalXML unmarshals the monetary value from an XML serialization
func (a *AmountValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	value := struct{ V string }{}
	d.DecodeElement(&value, &start)
	decoded, err := NewAmountValue(value.V)
	if err != nil {
		return err
	}
	*a = *decoded
	return nil
}

// MarshalXML marshals the monetary value for XML serialization
func (a *AmountValue) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(a.String(), start)
	return nil
}

func (t *Transaction) String() string {
	out := `TRANSACTION
  From:           %s\n
  To:             %s\n
  At:             %s\n
  Tags:           %s\n
  Amount:         %s\n
`
	return fmt.Sprintf(out, t.App.Entities[t.From], t.App.Entities[t.To], t.Datetime, strings.Join(t.Tags, ", "), t.Amount)
}

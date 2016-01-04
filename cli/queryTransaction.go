package cli

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/meisterluk/personal-accounting.go/types"
)

func queryTransaction(t *Transaction, app *Application) {
	reader := bufio.NewReader(os.Stdin)

	// print entities
	fmt.Println("Entities:")
	var min, max int64 = math.MaxInt64, 0
	for i, e := range app.Entities {
		if i > max {
			max = i
		}
		if i == -1 || i < min {
			min = i
		}
		fmt.Printf("  (%d) - %s\n", i, e)
	}
	fmt.Println()

	// query user
	fmt.Println("Money is transferred")
	// source entity
	fmt.Printf(" from entity [%d-%d]: ", min, max)
	from, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Sprintf("Error while reading from stdin: %s", err.Error()))
	}
	srcEntity, err := retrieveEntity(from, app)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf(" to entity [%d-%d]: ", min, max)
	to, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Sprintf("Error while reading from stdin: %s", err.Error()))
	}
	// destination entity
	destEntity, err := retrieveEntity(to, app)
	if err != nil {
		panic(err.Error())
	}
	// timestamp
	now := time.Now().Format(time.RFC3339)
	fmt.Printf(" at date [%s]: ", now)
	timestamp, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Sprintf("Error while reading from stdin: %s", err.Error()))
	}
	var when time.Time
	if timestamp == "\n" {
		when = time.Now()
	} else {
		when, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			panic(fmt.Sprintf("Could not parse date: %s", err.Error()))
		}
	}
	// Amount
	fmt.Printf(" value: ")
	amt, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Sprintf("Error while reading from stdin: %s", err.Error()))
	}
	amount, err := types.NewAmountValue(amt)
	if err != nil {
		panic(err.Error())
	}

	// tags
	fmt.Println()
	fmt.Println("Tags:")
	for i, tag := range app.Tags {
		fmt.Printf("  (%d) - %s\n", i, tag)
	}
	fmt.Println()
	fmt.Printf("With tags [space-separated]: ")
	tagIDs, err := reader.ReadString('\n')
	var tags []string
	for _, tagID := range strings.Split(tagIDs, " ") {
		id, err := strconv.ParseInt(strings.TrimSpace(tagID), 10, 64)
		if err != nil {
			continue
		}
		tagname := app.Tags[id]
		tagnames := strings.Split(tagname, ";")
		for _, tag := range tagnames {
			if !contains(tags, strings.TrimSpace(tag)) {
				tags = append(tags, strings.TrimSpace(tag))
			}
		}
	}

	fmt.Println()
	fmt.Println("TRANSACTION")
	fmt.Printf("  From:           %s\n", app.Entities[srcEntity])
	fmt.Printf("  To:             %s\n", app.Entities[destEntity])
	fmt.Printf("  At:             %s\n", when.Format(time.RFC3339))
	fmt.Printf("  Tags:           %s\n", strings.Join(tags, ", "))
	fmt.Printf("  Amount:         %s\n", amount)
	fmt.Println()

	t.Amount = amount
	t.Datetime = TransactionTimestamp(when)
	t.From = srcEntity
	t.To = destEntity
	t.Tags = tags
}

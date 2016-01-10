package cli

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/meisterluk/personal-accounting-go/types"
	"github.com/meisterluk/personal-accounting-go/utils"
)

func retrieveEntity(entityID string, app *types.Application) (int64, error) {
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

func queryTransaction(t *types.Transaction, app *types.Application) {
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
		when, err = time.Parse(time.RFC3339, strings.TrimSpace(timestamp))
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
			if !utils.ContainsString(tags, strings.TrimSpace(tag)) {
				tags = append(tags, strings.TrimSpace(tag))
			}
		}
	}

	fmt.Println(t.String())

	t.Amount = *amount
	t.Datetime = types.TransactionTimestamp{App: app, DateTime: when}
	t.From = srcEntity
	t.To = destEntity
	t.Tags = tags
}

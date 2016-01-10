package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/meisterluk/personal-accounting-go/types"
)

func printMenu() {
	fmt.Println(" (n) append new transaction")
	fmt.Println(" (l) list transactions")
	fmt.Println(" (q) quit / exit")
	fmt.Println()
}

// MainLoop for the CLI REPL
func MainLoop(app *types.Application, s *types.Storage) error {
	reader := bufio.NewReader(os.Stdin)

	for true {
		printMenu()
		inp, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input := strings.TrimSpace(inp)
		if input == "q" {
			fmt.Println("Bye!")
			break
		} else if input == "n" {
			t := types.Transaction{}
			queryTransaction(&t, app)
			err := (*s).Add(app, &t)
			if err != nil {
				return err
			}
			fmt.Println("Successfully done.")
		} else if input == "l" {
			tchan := make(chan types.Transaction)
			go func() {
				(*s).List(app, tchan)
				close(tchan)
			}()
			for {
				t, more := <-tchan
				if more {
					listTransaction(app, &t)
				} else {
					break
				}
			}
		}
	}

	return nil
}

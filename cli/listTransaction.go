package cli

import (
	"fmt"

	"github.com/meisterluk/personal-accounting-go/types"
)

func listTransaction(app *types.Application, t *types.Transaction) error {
	fmt.Println(t.String())
	return nil
}

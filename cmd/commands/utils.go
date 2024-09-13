package commands

import (
	"fmt"
	"os"
)

func ExitOnError(err error) {
	fmt.Printf("immortal error: %s\n", err.Error()) //nolint
	os.Exit(1)
}

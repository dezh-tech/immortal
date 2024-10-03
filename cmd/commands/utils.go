package commands

import (
	"log"
	"os"
)

func ExitOnError(err error) {
	log.Printf("immortal error: %s\n", err.Error()) //nolint
	os.Exit(1)
}

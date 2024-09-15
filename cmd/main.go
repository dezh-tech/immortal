package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/dezh-tech/immortal"
	"github.com/dezh-tech/immortal/cmd/commands"
)

func main() {
	if len(os.Args) < 2 {
		commands.ExitOnError(errors.New("at least 1 arguments expected.\nuse help command for more information"))
	}

	switch os.Args[1] {
	case "run":
		commands.HandleRun(os.Args)

	case "help":
		commands.HandleHelp(os.Args)
		os.Exit(0)

	case "version":
		fmt.Println(immortal.StringVersion()) //nolint
		os.Exit(0)

	default:
		commands.HandleHelp(os.Args)
	}
}

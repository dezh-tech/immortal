package commands

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/relay"
)

func HandleRun(args []string) {
	if len(args) < 3 {
		ExitOnError(errors.New("at least 1 arguments expected\nuse help command for more information"))
	}

	cfg, err := config.LoadfromFile(args[2])
	if err != nil {
		ExitOnError(err)
	}

	r := relay.NewRelay(*cfg)
	if err := r.Start(); err != nil {
		ExitOnError(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigChan

	if err := r.Stop(); err != nil {
		ExitOnError(err)
	}
}

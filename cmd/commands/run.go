package commands

import (
	"errors"
	"log"
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

	log.Println("loading config...")

	cfg, err := config.Load(args[2])
	if err != nil {
		ExitOnError(err)
	}

	r, err := relay.New(cfg)
	if err != nil {
		ExitOnError(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	errCh := r.Start()

	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %s\nInitiating graceful shutdown...\n", sig.String()) //nolint
		if err := r.Stop(); err != nil {
			ExitOnError(err)
		}

	case err := <-errCh:
		log.Printf("Unexpected error: %v\nInitiating shutdown due to the error...\n", err) //nolint
		if err := r.Stop(); err != nil {
			ExitOnError(err)
		}
	}
}

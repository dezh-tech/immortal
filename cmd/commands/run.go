package commands

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/dezh-tech/immortal/cmd/relay"
	"github.com/dezh-tech/immortal/config"
	"github.com/dezh-tech/immortal/pkg/logger"
)

func HandleRun(args []string) {
	if len(args) < 3 {
		ExitOnError(errors.New("at least 1 arguments expected\nuse help command for more information"))
	}

	cfg, err := config.Load(args[2])
	if err != nil {
		ExitOnError(err)
	}

	logger.InitGlobalLogger(&cfg.Logger)

	r, err := relay.New(cfg)
	if err != nil {
		ExitOnError(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	errCh := r.Start()

	select {
	case sig := <-sigChan:
		logger.Info("Received signal: Initiating graceful shutdown", "signal", sig.String())
		if err := r.Stop(); err != nil {
			ExitOnError(err)
		}

	case err := <-errCh:
		logger.Info("Unexpected error: Initiating shutdown due to the error", "err", err)
		if err := r.Stop(); err != nil {
			ExitOnError(err)
		}
	}
}

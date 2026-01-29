package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/pixil98/go-log"
)

type Validator interface {
	Validate() error
}

type WorkerBuilder func(any) (WorkerList, error)

type App struct {
	workers WorkerList
	logger  logrus.FieldLogger
}

func NewApp(config Validator, wb WorkerBuilder) (*App, error) {
	opts, err := newCmdLineOpts(os.Args[0])
	if err != nil {
		return nil, fmt.Errorf("configuring cmdline opts: %w", err)
	}

	err = opts.Parse(os.Args)
	if err != nil {
		return nil, fmt.Errorf("parsing cmdline opts: %w", err)
	}

	logger, err := opts.Logger()
	if err != nil {
		return nil, fmt.Errorf("configuring logger: %w", err)
	}

	err = opts.Config(config)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	workers, err := wb(config)
	if err != nil {
		return nil, fmt.Errorf("building workers: %w", err)
	}

	logger.Info("application initialized")

	return &App{
		workers: workers,
		logger:  logger,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	if a.workers == nil {
		return nil
	}

	ctx = log.SetLogger(ctx, a.logger)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	return a.workers.Start(ctx)
}

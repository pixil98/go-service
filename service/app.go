package service

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/pixil98/go-log/log"
)

type Valitator interface {
	Validate() error
}

type WorkerBuilder func(interface{}) (WorkerList, error)

type App struct {
	workers WorkerList
	logger  logrus.FieldLogger
}

func NewApp(config Valitator, wb WorkerBuilder) (*App, error) {
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
	//TODO: Support ctrl-c somehow

	return a.workers.Start(ctx)
}

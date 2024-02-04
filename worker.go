package service

import (
	"context"
	"sync"

	"github.com/pixil98/go-errors"
	"github.com/pixil98/go-log"
)

type Worker interface {
	Start(ctx context.Context) error
}

type WorkerList map[string]Worker

func (wl *WorkerList) Start(ctx context.Context) error {
	logger := log.GetLogger(ctx)

	// Cancel all workers if any one exits
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errs := errors.NewErrorList()
	for name, worker := range *wl {
		wg.Add(1)

		go func(name string, w Worker) {
			defer cancel()
			defer wg.Done()

			logger.Infof("starting %s", name)

			err := w.Start(ctx)
			if err != nil {
				errs.Add(err)
			}

			logger.Infof("exiting %s", name)
		}(name, worker)
	}

	wg.Wait()

	return errs.Err()
}

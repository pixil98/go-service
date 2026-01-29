package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/pixil98/go-errors"
)

type Worker interface {
	Start(ctx context.Context) error
}

type WorkerList map[string]Worker

func (wl *WorkerList) Start(ctx context.Context) error {
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

			slog.InfoContext(ctx, "starting worker", slog.String("name", name))

			err := w.Start(ctx)
			if err != nil {
				errs.Add(err)
			}

			slog.InfoContext(ctx, "exiting worker", slog.String("name", name))
		}(name, worker)
	}

	wg.Wait()

	return errs.Err()
}

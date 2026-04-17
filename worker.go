package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

type Worker interface {
	Start(ctx context.Context) error
}

type WorkerList map[string]Worker

func (wl *WorkerList) Start(ctx context.Context) error {
	// Cancel all workers if any one exits
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)
	for name, worker := range *wl {
		wg.Add(1)

		go func(name string, w Worker) {
			defer cancel()
			defer wg.Done()

			slog.InfoContext(ctx, "starting worker", slog.String("name", name))

			err := w.Start(ctx)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}

			slog.InfoContext(ctx, "exiting worker", slog.String("name", name))
		}(name, worker)
	}

	wg.Wait()

	return errors.Join(errs...)
}

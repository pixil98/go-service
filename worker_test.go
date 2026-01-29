package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"
)

type mockWorker struct {
	startFunc func(ctx context.Context) error
}

func (m *mockWorker) Start(ctx context.Context) error {
	if m.startFunc != nil {
		return m.startFunc(ctx)
	}
	return nil
}

func init() {
	// Suppress log output during tests
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestWorkerList_Start(t *testing.T) {
	tests := map[string]struct {
		workers WorkerList
		wantErr bool
	}{
		"empty worker list": {
			workers: WorkerList{},
			wantErr: false,
		},
		"single worker success": {
			workers: WorkerList{
				"worker1": &mockWorker{},
			},
			wantErr: false,
		},
		"single worker error": {
			workers: WorkerList{
				"worker1": &mockWorker{
					startFunc: func(ctx context.Context) error {
						return errors.New("worker failed")
					},
				},
			},
			wantErr: true,
		},
		"multiple workers all succeed": {
			workers: WorkerList{
				"worker1": &mockWorker{},
				"worker2": &mockWorker{},
			},
			wantErr: false,
		},
		"multiple workers one fails": {
			workers: WorkerList{
				"worker1": &mockWorker{},
				"worker2": &mockWorker{
					startFunc: func(ctx context.Context) error {
						return errors.New("worker2 failed")
					},
				},
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.workers.Start(t.Context())
			if (err != nil) != tt.wantErr {
				t.Errorf("WorkerList.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkerList_ContextCancellation(t *testing.T) {
	var cancelled atomic.Bool

	workers := WorkerList{
		"blocking": &mockWorker{
			startFunc: func(ctx context.Context) error {
				<-ctx.Done()
				cancelled.Store(true)
				return ctx.Err()
			},
		},
	}

	ctx, cancel := context.WithCancel(t.Context())

	done := make(chan error)
	go func() {
		done <- workers.Start(ctx)
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case <-done:
		if !cancelled.Load() {
			t.Error("worker did not receive cancellation")
		}
	case <-time.After(time.Second):
		t.Error("workers did not stop after context cancellation")
	}
}

func TestWorkerList_OneExitCancelsOthers(t *testing.T) {
	var worker2Cancelled atomic.Bool

	workers := WorkerList{
		"quick": &mockWorker{
			startFunc: func(ctx context.Context) error {
				return nil
			},
		},
		"blocking": &mockWorker{
			startFunc: func(ctx context.Context) error {
				<-ctx.Done()
				worker2Cancelled.Store(true)
				return nil
			},
		},
	}

	done := make(chan error)
	go func() {
		done <- workers.Start(t.Context())
	}()

	select {
	case <-done:
		if !worker2Cancelled.Load() {
			t.Error("blocking worker was not cancelled when quick worker exited")
		}
	case <-time.After(time.Second):
		t.Error("workers did not complete in time")
	}
}

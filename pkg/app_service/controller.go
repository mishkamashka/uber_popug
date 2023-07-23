package app_service

import (
	"context"
	"errors"
	"sync"
)

type NestedService interface {
	OnStart(ctx context.Context) error
	OnStop()
}

type Controller struct {
	service NestedService

	onceStart sync.Once
	onceStop  sync.Once

	done     chan struct{}
	commands chan<- Command
}

func New(service NestedService) *Controller {
	return &Controller{
		service: service,
		done:    make(chan struct{}),
	}
}

func (ctrl *Controller) SetCommandsChan(commands chan<- Command) {
	ctrl.commands = commands
}

func (ctrl *Controller) Start(ctx context.Context) (err error) {
	started := false
	ctrl.onceStart.Do(func() {
		err = ctrl.service.OnStart(ctx)
		if err == nil {
			go func() {
				select {
				case <-ctrl.done:
					break
				case <-ctx.Done():
					_ = ctrl.Stop()
				}
			}()
		}

		started = true
	})

	if !started {
		return errors.New("already started")
	}
	return err
}

func (ctrl *Controller) Stop() error {
	stopped := false
	ctrl.onceStop.Do(func() {
		// Wait for service logic
		ctrl.service.OnStop()

		close(ctrl.done)
		stopped = true
	})

	if !stopped {
		return errors.New("already called")
	}
	return nil
}

func (ctrl *Controller) Done() <-chan struct{} {
	return ctrl.done
}

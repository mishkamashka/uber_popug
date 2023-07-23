package app_service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Service interface {
	Start(ctx context.Context) error
	Done() <-chan struct{}
	SetCommandsChan(chan<- Command)
}

type launcher struct {
	main            Service
	auxiliary       []Service
	cancelAuxiliary context.CancelFunc
	commands        chan Command
	readiness       uint32
	allReady        uint32
}

func NewLauncher(main Service, auxiliary ...Service) *launcher {
	l := launcher{
		main:      main,
		auxiliary: auxiliary,
	}
	return &l
}

func (l *launcher) IsReady() bool {
	return l.readiness == l.allReady
}

// Run runs service as app
func (l *launcher) Run(signals ...os.Signal) error {
	l.commands = make(chan Command, len(l.auxiliary))
	defer close(l.commands)

	ctx, cancel := context.WithCancel(context.Background())
	err := l.onStart(ctx)
	if err != nil {
		cancel()
		return err
	}

	err = l.waitForSignal(signals)

	cancel()
	l.onStop()

	log.Println("bye 👋")
	return err
}

func (l *launcher) onStart(ctx context.Context) error {
	auxiliaryCtx, cancelAuxiliary := context.WithCancel(context.Background())
	l.cancelAuxiliary = cancelAuxiliary
	for i, auxiliaryService := range l.auxiliary {
		l.allReady |= uint32(1) << (i + 1) // `0` reserved by main service
		auxiliaryService.SetCommandsChan(l.commands)

		if err := auxiliaryService.Start(auxiliaryCtx); err != nil {
			return fmt.Errorf("auxiliary service#%d start error: %w", i, err)
		}
	}

	l.allReady |= uint32(1)
	l.readiness = l.allReady
	l.main.SetCommandsChan(l.commands)
	if err := l.main.Start(ctx); err != nil {
		return fmt.Errorf("main service start error: %w", err)
	}

	return nil
}

func (l *launcher) onStop() {
	<-l.main.Done()
	l.cancelAuxiliary()
	for _, auxiliaryService := range l.auxiliary {
		<-auxiliaryService.Done()
	}
}

var defaultSignals = []os.Signal{syscall.SIGTERM, syscall.SIGINT}

func (l *launcher) waitForSignal(signals []os.Signal) error {
	if len(signals) == 0 {
		signals = defaultSignals
	}

	sig := make(chan os.Signal, len(signals))
	signal.Notify(sig, signals...)
	for {
		select {
		case <-l.main.Done():
			log.Println("main is stopped")
			return nil
		case cmd := <-l.commands:
			stop, err := cmd(l)
			if err != nil {
				log.Println("error in service")
				return fmt.Errorf("service error: %w", err)
			}
			if stop {
				return nil
			}
		case s := <-sig:
			log.Printf("received signal: %s", s.String())
			return nil
		}
	}
}

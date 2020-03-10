package relogger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/fatih/color"
)

var (
	// ErrInvalidParams presents an invalid error.
	ErrInvalidParams = errors.New("[err] invalid params")
)

// ReLogger is a custom logger to support refresh file using an interval or signal.
type ReLogger struct {
	log.Logger
	lock            sync.Mutex
	filepath        string
	file            *os.File
	filemode        os.FileMode
	signals         []os.Signal
	refreshDuration time.Duration
	printableDebug  bool
	signalBackOff   *backoff.ExponentialBackOff
	intervalBackOff *backoff.ExponentialBackOff
}

// NewReLogger returns a logger.
func NewReLogger(filepath string, opts ...Option) (*ReLogger, error) {
	if filepath == "" {
		return nil, fmt.Errorf("[err] NewReLogger %w", ErrInvalidParams)
	}

	logger := &ReLogger{filepath: filepath,
		signalBackOff:   backoff.NewExponentialBackOff(),
		intervalBackOff: backoff.NewExponentialBackOff()}

	mergeOpts := []Option{
		WithFileMode(os.ModePerm),
		WithSignals([]os.Signal{syscall.SIGHUP}),
		WithRefreshDuration(1 * time.Hour),
	}

	mergeOpts = append(mergeOpts, opts...)
	for _, opt := range mergeOpts {
		opt.apply(logger)
	}

	file, err := os.OpenFile(logger.filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, logger.filemode)
	if err != nil {
		return nil, fmt.Errorf("[err] NewReLogger %w", err)
	}
	logger.file = file
	logger.SetOutput(file)

	go logger.signalRoutine()
	go logger.intervalRoutine()

	return logger, nil
}

func (l *ReLogger) signalRoutine() {
	if len(l.signals) == 0 {
		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, l.signals...)
	for {
		s := <-sig
		if l.printableDebug {
			color.Cyan("[info][%s][relogger] signal (%s)", time.Now().Format(time.RFC3339), s.String())
		}

		for {
			time.Sleep(l.signalBackOff.NextBackOff())
			if err := l.refresh(); err != nil {
				if l.printableDebug {
					color.Red("%s", err.Error())
				}
				continue
			}
			l.signalBackOff.Reset()
			break
		}
	}
}

func (l *ReLogger) intervalRoutine() {
	if l.refreshDuration <= 0 {
		return
	}

	for {
		select {
		case <-time.After(l.refreshDuration):
			if l.printableDebug {
				color.Cyan("[info][%s][relogger] interval", time.Now().Format(time.RFC3339))
			}

			for {
				time.Sleep(l.intervalBackOff.NextBackOff())
				if err := l.refresh(); err != nil {
					if l.printableDebug {
						color.Red("%s", err.Error())
					}
					continue
				}
				l.intervalBackOff.Reset()
				break
			}
		}
	}
}

func (l *ReLogger) refresh() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.printableDebug {
		color.Cyan("[info][%s][relogger] refresh", time.Now().Format(time.RFC3339))
	}

	file, err := os.OpenFile(l.filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, l.filemode)
	if err != nil {
		return fmt.Errorf("[err][%s][relogger] refresh %w", time.Now().Format(time.RFC3339), err)
	}
	l.SetOutput(file)

	if l.file != nil {
		if err := l.file.Close(); err != nil {
			color.Red("[err][%s][relogger] old file close %s", time.Now().Format(time.RFC3339), err.Error())
		}
	}
	l.file = file

	if l.printableDebug {
		color.Green("[success][%s][relogger] refresh", time.Now().Format(time.RFC3339))
	}
	return nil
}

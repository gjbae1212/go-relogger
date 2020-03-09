package relogger

import (
	"os"
	"time"
)

type Option interface {
	apply(l *ReLogger)
}

type OptionFunc func(l *ReLogger)

func (o OptionFunc) apply(l *ReLogger) { o(l) }

// WithFileMode returns the func which sets file mode.
func WithFileMode(mode os.FileMode) OptionFunc {
	return func(l *ReLogger) {
		l.filemode = mode
	}
}

// WithSignals returns the func which traps kill signal.
func WithSignals(signals []os.Signal) OptionFunc {
	return func(l *ReLogger) {
		l.signals = signals
	}
}

// WithRefreshDuration returns the func which sets an interval for refreshing logger.
func WithRefreshDuration(d time.Duration) OptionFunc {
	return func(l *ReLogger) {
		l.refreshDuration = d
	}
}

// WithPrintableDebug returns the func which sets bool to check whether debug or not.
func WithPrintableDebug(debug bool) OptionFunc {
	return func(l *ReLogger) {
		l.printableDebug = debug
	}
}

package relogger

import (
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithFileMode(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		input os.FileMode
	}{
		"success": {input: os.ModePerm},
	}

	for _, t := range tests {
		l := &ReLogger{}
		f := WithFileMode(t.input)
		f(l)
		assert.Equal(l.filemode, t.input)
	}
}

func TestWithSignals(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		input []os.Signal
	}{
		"success": {input: []os.Signal{syscall.SIGHUP, syscall.SIGUSR1}},
	}

	for _, t := range tests {
		l := &ReLogger{}
		f := WithSignals(t.input)
		f(l)
		assert.True(reflect.DeepEqual(l.signals, t.input))
	}
}

func TestWithRefresh(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		input time.Duration
	}{
		"success": {input: time.Second * 10},
	}

	for _, t := range tests {
		l := &ReLogger{}
		f := WithRefreshDuration(t.input)
		f(l)
		assert.Equal(l.refreshDuration, t.input)
	}

}

func TestWithPrintableDebug(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		input bool
	}{
		"success": {input: true},
	}

	for _, t := range tests {
		l := &ReLogger{}
		f := WithPrintableDebug(t.input)
		f(l)
		assert.Equal(l.printableDebug, t.input)
	}
}

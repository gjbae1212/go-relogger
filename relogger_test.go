package relogger

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewReLogger(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, ok := runtime.Caller(0)
	assert.True(ok)

	tests := map[string]struct {
		filepath string
		opts     []Option
		isErr    bool
	}{
		"fail": {isErr: true},
		"success": {
			filepath: filepath.Join(path.Dir(filename), "test.log"),
			opts: []Option{
				WithPrintableDebug(true),
			}},
	}

	for _, t := range tests {
		_, err := NewReLogger(t.filepath, t.opts...)
		assert.Equal(t.isErr, err != nil)
	}

}

func TestCheckRefresh(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, ok := runtime.Caller(0)
	assert.True(ok)
	filepath := filepath.Join(path.Dir(filename), "test.log")

	logger, err := NewReLogger(filepath,
		WithRefreshDuration(5 * time.Second),
		WithPrintableDebug(true))
	assert.NoError(err)

	for i := 0; i < 60; i++ {
		time.Sleep(1 * time.Second)
		logger.Printf("[%s] TestCheckRefresh ", time.Now().Format(time.RFC3339))

		if i % 11 == 0 {

		}
	}
}

func BenchmarkNewReLogger(b *testing.B) {
	_, filename, _, _ := runtime.Caller(0)
	filepath := filepath.Join(path.Dir(filename), "test.log")
	logger, _ := NewReLogger(filepath,
		WithRefreshDuration(1 * time.Second),
		WithPrintableDebug(true))
	for i := 0; i < b.N; i++ {
		logger.Printf("[%s] BenchmarkNewReLogger ", time.Now().Format(time.RFC3339))
	}
}

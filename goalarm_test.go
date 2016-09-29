package goalarm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	startTime := time.Now()
	job := func(ch chan interface{}) {
		ch <- "hello"
	}
	ch := In(2*time.Second, job)
	assert.Equal(t, "hello", <-ch)
	assert.InDelta(t, int(2*time.Second), int(time.Now().Sub(startTime)), float64(100*time.Millisecond))
}

func TestAt(t *testing.T) {
	startTime := time.Now()
	job := func(ch chan interface{}) {
		ch <- "hello"
	}
	ch := At(time.Now().Add(2*time.Second), job)
	assert.Equal(t, "hello", <-ch)
	assert.InDelta(t, int(2*time.Second), int(time.Now().Sub(startTime)), float64(100*time.Millisecond))
}

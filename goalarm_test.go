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
	dur := int(time.Now().Sub(startTime))
	assert.InDelta(t, int(2*time.Second), dur, float64(100*time.Millisecond))
}

func TestAt(t *testing.T) {
	startTime := time.Now()
	job := func(ch chan interface{}) {
		ch <- "hello"
	}
	ch := At(time.Now().Add(2*time.Second), job)
	assert.Equal(t, "hello", <-ch)
	dur := int(time.Now().Sub(startTime))
	assert.InDelta(t, int(2*time.Second), dur, float64(100*time.Millisecond))
}

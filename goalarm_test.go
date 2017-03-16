package goalarm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	startTime := time.Now()
	job := func(ctx context.Context) (interface{}, error) {
		return "hello", nil
	}
	ctx := context.Background()
	ch := In(ctx, 1*time.Second, job)
	assert.Equal(t, "hello", <-ch)
	dur := int(time.Since(startTime))
	assert.InDelta(t, int(1*time.Second), dur, float64(100*time.Millisecond))
}

func TestAt(t *testing.T) {
	startTime := time.Now()
	job := func(ctx context.Context) (interface{}, error) {
		return "hello", nil
	}
	ctx := context.Background()
	ch := At(ctx, time.Now().Add(1*time.Second), job)
	assert.Equal(t, "hello", <-ch)
	dur := int(time.Since(startTime))
	assert.InDelta(t, int(1*time.Second), dur, float64(100*time.Millisecond))
}

func TestEvery(t *testing.T) {
	job := func(ctx context.Context) (interface{}, error) {
		return "hello", nil
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := Every(ctx, 1*time.Second, job)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	count := 0
Loop:
	for {
		select {
		case in := <-ch:
			switch in.(type) {
			case string:
				assert.Equal(t, "hello", in)
				count++
			case error:
				if assert.Equal(t, context.Canceled, in) {
					break Loop
				} else {
					panic(in)
				}
			default:
				panic(in)
			}
		}
	}
	assert.Equal(t, 2, count)
}

func TestEvery_ErrorShouldStopExecution(t *testing.T) {
	err := errors.New("error")
	job := func(ctx context.Context) (interface{}, error) {
		return "", err
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := Every(ctx, 1*time.Second, job)
	go func() {
		time.Sleep(4 * time.Second)
		cancel()
	}()
	count := 0
Loop:
	for {
		select {
		case in := <-ch:
			switch in.(type) {
			case string:
				assert.Equal(t, "hello", in)
				count++
			case error:
				if assert.Equal(t, err, in) {
					break Loop
				} else {
					panic(in)
				}
			default:
				panic(in)
			}
		}
	}
	assert.Equal(t, 0, count)
}

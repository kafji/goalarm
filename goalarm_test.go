package goalarm

import (
	"context"
	"errors"
	"sync"
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

func TestIn_ShouldSendErrorWhenJobReturnError(t *testing.T) {
	job := func(ctx context.Context) (interface{}, error) {
		return nil, errors.New("error")
	}
	ctx := context.Background()
	ch := In(ctx, 1*time.Second, job)
	assert.Equal(t, errors.New("error"), <-ch)
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

func TestAt_ShouldSendErrorWhenJobReturnError(t *testing.T) {
	job := func(ctx context.Context) (interface{}, error) {
		return nil, errors.New("error")
	}
	ctx := context.Background()
	ch := At(ctx, time.Now().Add(1*time.Second), job)
	assert.Equal(t, errors.New("error"), <-ch)
}

func TestEvery(t *testing.T) {
	startTime := time.Now()
	job := func(ctx context.Context) (interface{}, error) {
		time.Sleep(1 * time.Second)
		return "hello", nil
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := Every(ctx, 2*time.Second, 0, job)

	var cancelWG sync.WaitGroup
	cancelWG.Add(3)
	go func() {
		cancelWG.Wait()
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
				cancelWG.Done()
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

	assert.Equal(t, 3, count)
	dur := int(time.Since(startTime))
	assert.InDelta(t, int(7*time.Second), dur, float64(100*time.Millisecond))
}

func TestEvery_ExecuteFirstJobWhenCalled(t *testing.T) {
	job := func(ctx context.Context) (interface{}, error) {
		return "hello", nil
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := Every(ctx, 1*time.Second, 0, job)

	go func() {
		time.Sleep(100 * time.Millisecond)
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
	assert.Equal(t, 1, count)
}

func TestEvery_WithFirstDelay(t *testing.T) {
	startTime := time.Now()
	job := func(ctx context.Context) (interface{}, error) {
		time.Sleep(1 * time.Second)
		return "hello", nil
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := Every(ctx, 2*time.Second, 2*time.Second, job)

	var cancelWG sync.WaitGroup
	cancelWG.Add(3)
	go func() {
		cancelWG.Wait()
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
				cancelWG.Done()
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
	assert.Equal(t, 3, count)
	dur := int(time.Since(startTime))
	assert.InDelta(t, int(9*time.Second), dur, float64(100*time.Millisecond))
}

func TestEvery_ErrorShouldStopExecution(t *testing.T) {
	err := errors.New("error")
	job := func(ctx context.Context) (interface{}, error) {
		return "", err
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ch := Every(ctx, 1*time.Second, 0, job)

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

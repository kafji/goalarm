package goalarm

import (
	"context"
	"errors"
	"log"
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
	dur := int(time.Now().Sub(startTime))
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
	dur := int(time.Now().Sub(startTime))
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
		time.Sleep(4 * time.Second)
		cancel()
	}()
	count := 0
Loop:
	for {
		select {
		case in := <-ch:
			if s, ok2 := in.(string); ok2 {
				log.Println(s)
				assert.Equal(t, "hello", s)
				count++
			} else if e, ok2 := in.(error); ok2 {
				if assert.Equal(t, context.Canceled, e) {
					break Loop
				} else {
					panic(e)
				}
			} else {
				panic(in)
			}
		}
	}
	assert.Equal(t, 3, count)
}

func TestEveryErrorShouldStopExecution(t *testing.T) {
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
			if s, ok2 := in.(string); ok2 {
				log.Println(s)
				assert.Equal(t, "hello", s)
				count++
			} else if e, ok2 := in.(error); ok2 {
				if assert.Equal(t, err, e) {
					break Loop
				} else {
					panic(e)
				}
			} else {
				panic(in)
			}
		}
	}
	assert.Equal(t, 0, count)
}

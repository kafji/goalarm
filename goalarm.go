package goalarm

import (
	"context"
	"sync"
	"time"
)

type Job func(context.Context) (interface{}, error)

func In(ctx context.Context, d time.Duration, job Job) chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		time.Sleep(d)
		if r, err := job(ctx); err != nil {
			ch <- err
		} else {
			ch <- r
		}
	}()
	return ch
}

func At(ctx context.Context, t time.Time, job Job) chan interface{} {
	d := t.Sub(time.Now())
	return In(ctx, d, job)
}

func Every(ctx context.Context, delay time.Duration, job Job) chan interface{} {
	var wg sync.WaitGroup
	ch := make(chan interface{})
	executionCh := make(chan interface{})
	var run func()
	run = func() {
		defer wg.Done()
		resultCh := In(ctx, delay, job)
		executionCh <- <-resultCh
	}
	go func() {
		for {
			select {
			case result := <-executionCh:
				ch <- result
				wg.Add(1)
				go run()
			case <-ctx.Done():
				ch <- ctx.Err()
				wg.Wait()
				close(ch)
				close(executionCh)
			}
		}
	}()
	wg.Add(1)
	go run()
	return ch
}

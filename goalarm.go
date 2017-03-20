package goalarm

import (
	"context"
	"sync"
	"time"
)

type Job func(context.Context) (interface{}, error)

func In(ctx context.Context, d time.Duration, job Job) chan interface{} {
	outbound := make(chan interface{}, 1)
	go func() {
		defer close(outbound)
		time.Sleep(d)
		if r, err := job(ctx); err != nil {
			outbound <- err
		} else {
			outbound <- r
		}
	}()
	return outbound
}

func At(ctx context.Context, t time.Time, job Job) chan interface{} {
	d := t.Sub(time.Now())
	return In(ctx, d, job)
}

func Every(ctx context.Context, delay, firstDelay time.Duration, job Job) chan interface{} {
	var wg sync.WaitGroup

	outbound := make(chan interface{}, 1)

	executeCh := make(chan interface{})
	scheduleCh := make(chan interface{})

	// scheduler
	go func() {
		for {
			select {
			case <-scheduleCh:
				time.Sleep(delay)
				executeCh <- nil
			}
		}
	}()

	execute := func() {
		wg.Add(1)
		if r, err := job(ctx); err != nil {
			outbound <- err
		} else {
			outbound <- r
		}
		wg.Done()
	}

	// executor
	go func() {
		for {
			select {
			case <-executeCh:
				execute()
				scheduleCh <- nil
			case <-ctx.Done():
				outbound <- ctx.Err()
				wg.Wait()
				close(outbound)
				return
			}
		}
	}()

	// start
	go func() {
		time.Sleep(firstDelay)
		executeCh <- nil
	}()

	return outbound
}

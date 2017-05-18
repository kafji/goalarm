package goalarm

import (
	"context"
	"sync"
	"time"
)

// Job is function that will get executed.
// Return error to stop execution.
// Result will be propragate to channel returned by scheduler functions if err is nil.
type Job func(ctx context.Context) (result interface{}, err error)

// In is a scheduler function that will execute job after specified duration has passed.
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

// At is a scheduler function that will execute job at specified time.
func At(ctx context.Context, t time.Time, job Job) chan interface{} {
	d := t.Sub(time.Now())
	return In(ctx, d, job)
}

// Every is a scheduler function that will execute job periodically after specified delay has passed.
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

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

func Every(ctx context.Context, d time.Duration, job Job) chan interface{} {
	var wg sync.WaitGroup
	inbound := make(chan interface{})
	outbound := make(chan interface{}, 1)
	execute := func() {
		defer wg.Done()
		wg.Add(1)
		inbound <- <-In(ctx, d, job)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				outbound <- ctx.Err()
				wg.Wait()
				close(outbound)
				return
			case result := <-inbound:
				outbound <- result
				go execute()
			}
		}
	}()
	go execute()
	return outbound
}

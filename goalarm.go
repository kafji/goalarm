package goalarm

import "time"

type Job func(chan interface{})

func In(d time.Duration, job Job) chan interface{} {
	ch := make(chan interface{})
	go func() {
		time.Sleep(d)
		job(ch)
	}()
	return ch
}

func At(t time.Time, job Job) chan interface{} {
	d := t.Sub(time.Now())
	return In(d, job)
}

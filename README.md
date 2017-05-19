# goalarm

[![Build Status](https://travis-ci.org/kafji/goalarm.svg?branch=master)](https://travis-ci.org/kafji/goalarm)
[![codecov](https://codecov.io/gh/kafji/goalarm/branch/master/graph/badge.svg)](https://codecov.io/gh/kafji/goalarm)

Run _job_ periodically or at specific time.
Job is a function that accept `context.Context` and return `interface{}` and `error`.
Job will be executed in different goroutine than the caller.
Cancel periodically executed job using cancel function returned from `context.WithCancel(context.Background())`. For more use case take a look at its [tests](https://github.com/kafji/goalarm/blob/master/goalarm_test.go).

## Installation
```
go get -u github.com/kafji/goalarm
```

## Example
```go
ctx := context.Background()

// Print "hello" in 10 seconds.
goalarm.In(ctx, 10*time.Second, func(ctx context.Context) (interface{}, error) {
	fmt.Println("hello")
	return nil, nil
})

// Print "hello" every 10 seconds.
goalarm.Every(ctx, 10*time.Second, 0, func(ctx context.Context) (interface{}, error) {
	fmt.Println("hello")
	return nil, nil
})
```

## Development

### Run test
```
go test $(go list ./... | grep -v /vendor/)
```

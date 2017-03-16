# goalarm

## Installation
```
go get -u github.com/kafji/goalarm
```

## Example
```go
// Print "hello" in 10 seconds.
In(10*time.Second, func(ch chan interface{}) {
	defer close(ch)
	fmt.Println("hello")
})
```


## Development

### Run test
```
go test $(go list ./... | grep -v /vendor/)
```

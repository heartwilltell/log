# `log` - Simple wrapper around standard log package

- 😚 Simple API
- 👌 Zero dependencies
- 😮‍💨 No global logger
- 👏 No structured logging bullshit

## Installation
```bash
go get github.com/heartwilltell/log
```

## Usage

👇 Creates logger with `info` level. 
```go
logger := log.New()
```
###

👇 Creates nop logger which implements `log.Logger` interface.
```go
logger := log.NewNopLog()
```
💡 _Useful for tests or places where logger should be disabled by default_
###

👇 Creates logger with `debug` level.
```go
logger := log.New(log.WithLevel(log.DBG))
```
###

👇 Parses string to level and creates logger with `warning` level.
```go
level, levelErr := log.ParseLevel("warning")
if levelErr != nil {
	// handle error here
}

logger := log.New(log.WithLevel(level))
```
###

👇 Creates logger with different `io.Writer`.
```go
var buf []byte
w := bytes.NewBuffer(buf)

logger := log.New(log.WithWriter(w))
```
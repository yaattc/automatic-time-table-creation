# timetype ![Go](https://github.com/Semior001/timetype/workflows/Go/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/Semior001/timetype/badge.svg?branch=master)](https://coveralls.io/github/Semior001/timetype?branch=master) [![go report card](https://goreportcard.com/badge/github.com/semior001/timetype)](https://goreportcard.com/report/github.com/semior001/timetype) [![PkgGoDev](https://pkg.go.dev/badge/github.com/Semior001/timetype)](https://pkg.go.dev/github.com/Semior001/timetype)
Package adds some time types for easier work, serialize and deserialize them and some helper functions. Types satisfy the `fmt.GoStringer` and `fmt.Stringer` interfaces for easier debugging and `sql.Scanner` and `sql.Valuer` to allow to use this types with the SQL drivers.   

## `timetype.Clock`

The type implements `sql.Scanner` and `json.Unmarshaler` and tries to read the time value in two formats: ISO8601 for times without date 
and ISO8601 with micro precision without date.

```go
// Clock is a wrapper for time.time to allow parsing datetime stamp with time only in
// ISO 8601 format, like "15:04:05"
type Clock time.Time
```

```go
// NewClock returns the Clock in the given location with given hours, minutes and secs
func NewClock(h, m, s int, loc *time.Location) Clock
```

```go
// NewUTCClock returns new clock with given hours, minutes and seconds in the UTC location
func NewUTCClock(h, m, s int) Clock 
```

## `timetype.Duration`

```go
// Duration is a wrapper of time.Duration, that allows to marshal and unmarshal time in RFC3339 format
type Duration time.Duration
``` 

## Helpers

```go
// ParseWeekday parses a weekday from a string and, if it's
// can't be parsed, returns
func ParseWeekday(s string) (time.Weekday, error)
```

## Errors

```go
// Parsing errors
var (
    ErrInvalidClock    = errors.New("timetype: invalid clock")
    ErrInvalidDuration = errors.New("timetype: invalid duration")
    ErrInvalidWeekday  = errors.New("timetype: invalid weekday")
    ErrUnknownFormat   = errors.New("timetype: unknown format")
)
```

## Time formats
```go
// Templates to parse clocks
const (
	ISO8601Clock      = "15:04:05"
	ISO8601ClockMicro = "15:04:05.000000"
)
```

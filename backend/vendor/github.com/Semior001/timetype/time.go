package timetype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Parsing errors
var (
	ErrInvalidClock    = errors.New("timetype: invalid clock")
	ErrInvalidDuration = errors.New("timetype: invalid duration")
)

// Templates to parse clocks
const (
	ISO8601Clock      = "15:04:05"
	ISO8601ClockMicro = "15:04:05.000000"
)

// Clock is a wrapper for time.time to allow parsing datetime stamp with time only in
// ISO 8601 format, like "15:04:05"
type Clock time.Time

// NewClock returns the Clock in the given location with given hours, minutes and secs
func NewClock(h, m, s, ns int, loc *time.Location) Clock {
	return Clock(time.Date(0, time.January, 1, h, m, s, ns, loc))
}

// NewUTCClock returns new clock with given hours, minutes and seconds in the UTC location
func NewUTCClock(h, m, s, ns int) Clock {
	return NewClock(h, m, s, ns, time.UTC)
}

// MarshalJSON marshals time into time
func (h Clock) MarshalJSON() ([]byte, error) {
	res, err := json.Marshal(time.Time(h).Format(ISO8601ClockMicro))
	return res, wrapExternalErr(err)
}

// String implements fmt.Stringer to print and log Clock properly
func (h Clock) String() string {
	t := time.Time(h)
	return fmt.Sprintf("%02d:%02d:%02d %s", t.Hour(), t.Minute(), t.Second(), t.Location())
}

// GoString implements fmt.GoStringer to use Clock in %#v formats
func (h Clock) GoString() string {
	t := time.Time(h)
	return fmt.Sprintf("timetype.NewClock(%d, %d, %d, %s)", t.Hour(), t.Minute(), t.Second(), t.Location())
}

// UnmarshalJSON converts time to ISO 8601 representation
func (h *Clock) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return wrapExternalErr(err)
	}
	val, ok := v.(string)
	if !ok {
		return ErrInvalidClock
	}
	t, err := time.Parse(ISO8601ClockMicro, val)
	if err != nil {
		return wrapExternalErr(err)
	}
	*h = Clock(t)
	return nil
}

// Scan the given SQL value as Clock
func (h *Clock) Scan(src interface{}) (err error) {
	switch v := src.(type) {
	case nil:
		*h = Clock{}
	case time.Time:
		*h = Clock(v)
	case string:
		t, err := time.Parse(ISO8601ClockMicro, v)
		if err != nil {
			return wrapExternalErr(err)
		}
		*h = Clock(t)
	case []byte:
		t, err := time.Parse(ISO8601ClockMicro, string(v))
		if err != nil {
			return wrapExternalErr(err)
		}
		*h = Clock(t)
	default:
		return ErrInvalidClock
	}

	return err
}

// Value returns the SQL value of the given Clock
func (h Clock) Value() (driver.Value, error) {
	return time.Time(h).Format(ISO8601ClockMicro), nil
}

// Duration is a wrapper of time.Duration, that allows to marshal and unmarshal time in RFC3339 format
type Duration time.Duration

// MarshalJSON simply marshals duration into nanoseconds
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON converts time duration from RFC3339 format into time.Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return wrapExternalErr(err)
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return wrapExternalErr(err)
		}
		*d = Duration(tmp)
		return nil
	default:
		return ErrInvalidDuration
	}
}

// Scan the given SQL value as Duration
func (d *Duration) Scan(src interface{}) (err error) {
	switch v := src.(type) {
	case nil:
		*d = 0
	case time.Duration:
		*d = Duration(v)
	case float64:
		*d = Duration(v)
	case int64:
		*d = Duration(v)
	case string:
		err = wrapExternalErr(d.UnmarshalJSON([]byte(v)))
	case []byte:
		err = wrapExternalErr(d.UnmarshalJSON(v))
	default:
		return ErrInvalidDuration
	}

	return err
}

// Value returns the SQL value of the given Duration
func (d Duration) Value() (driver.Value, error) {
	return int64(d), nil
}

// errExternal wraps an error come outside this package (e.g. from time.ParseDuration).
// It allows to detect the external error inside tests by asserting the type of an error.
type errExternal struct {
	error // wrapped error
}

// Error returns the error string of the wrapped error
func (e *errExternal) Error() string {
	return e.error.Error()
}

func wrapExternalErr(e error) error {
	if e == nil {
		return nil
	}
	return &errExternal{error: e}
}

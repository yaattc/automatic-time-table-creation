package rest

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-pkgz/rest/logger"
)

// Validatable describes structs needed to validate
type Validatable interface {
	Validate() error
}

// Optimizable describes structs which might be optimized before passing to service methods
type Optimizable interface {
	Optimize()
}

// ISO8601Clock describes time layout in ISO 8601 standard
const ISO8601Clock = "15:04:05"

// JSON type represents every JSON-serializable map, just for convenience
type JSON map[string]interface{}

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
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// Clock is a wrapper for time.time to allow parsing datetime stamp with time only in
// ISO 8601 format, like "15:04:05"
type Clock time.Time

// MarshalJSON marshals time into time
func (h Clock) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(h).Format(ISO8601Clock))
}

// UnmarshalJSON converts time to ISO 8601 representation
func (h *Clock) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	val, ok := v.(string)
	if !ok {
		return errors.New("invalid clock")
	}
	t, err := time.Parse(ISO8601Clock, val)
	if err != nil {
		return err
	}
	*h = Clock(t)
	return nil
}

// stdLogger to implement logger.Backend function to allow to log in middleware
type stdLogger struct{}

// StdLogger returns standard backend logger
func StdLogger() logger.Backend {
	return stdLogger{}
}

// Logf simply calls the standard logger
func (s stdLogger) Logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

package timetype

import (
	"errors"
	"time"
)

// ErrInvalidWeekday if the weekday string cannot be parsed into time.Weekday
var ErrInvalidWeekday = errors.New("timetype: invalid weekday")

// weekdays to string names mapping
var weekdays = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

// ParseWeekday parses a weekday from a string and, if it's
// can't be parsed, returns
func ParseWeekday(s string) (time.Weekday, error) {
	if wd, ok := weekdays[s]; ok {
		return wd, nil
	}
	return 0, ErrInvalidWeekday
}

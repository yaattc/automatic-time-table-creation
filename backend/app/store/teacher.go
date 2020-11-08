package store

import (
	"time"

	"github.com/Semior001/timetype"
)

// Teacher describes a basic teacher with its own name and surname
type Teacher struct {
	Preferences TeacherPreferences `json:"preferences,omitempty"`
	TeacherDetails
}

// PrepareUntrusted sets zero values for all fields that are immutable for user
func (t *Teacher) PrepareUntrusted() {
	t.ID = ""
	t.Preferences = TeacherPreferences{}
}

// TeacherDetails describes a data that relates to one particular teacher
// to exclude the recursion problems
type TeacherDetails struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`

	Email string `json:"email"`

	Degree string `json:"degree"`
	About  string `json:"about"`
}

// TeacherPreferences describes teacher's preferences in generating the schedule.
// When setting the teacher preferences, in Staff all fields will be ignored except the IDs
type TeacherPreferences struct {
	TimeSlots []TimeSlot       `json:"time_slots"` // preferable time slots for teaching
	Staff     []TeacherDetails `json:"staff"`      // preferable teaching staff
	Locations []Location       `json:"rooms"`      // preferable rooms for teaching
}

// Empty checks whether the preferences are empty or not
func (p TeacherPreferences) Empty() bool {
	return len(p.TimeSlots) < 1 && len(p.Staff) < 1 && len(p.Locations) < 1
}

// TimeSlot describes a particular period of time in a week
type TimeSlot struct {
	Weekday  time.Weekday      `json:"weekday"`  // a weekday of time slot
	Start    timetype.Clock    `json:"start"`    // start time of time slot
	Duration timetype.Duration `json:"duration"` // duration of a time slot
	Location Location          `json:"location"` // an optional location field, empty means "any"
}

package store

import (
	"time"

	"github.com/Semior001/timetype"
)

// Location describes a room or auditory where the Class is held
type Location string

// EducationalProgram describes a study level of students
type EducationalProgram string

// Basic educational levels
const (
	Bachelor EducationalProgram = "bachelor"
	Master   EducationalProgram = "master"
)

// CourseFormat describes a parts (couples, lessons) of the Course
type CourseFormat string

// Basic course formats
const (
	Lecture  CourseFormat = "lecture"
	Tutorial CourseFormat = "tutorial"
	Lab      CourseFormat = "lab"
)

// Course describes a basic semester course, e.g. "Operational systems"
type Course struct {
	ID                 string             `json:"id"`                  // a hash derived from all others fields
	Name               string             `json:"name"`                // the name of the course
	Program            EducationalProgram `json:"program"`             // bachelor, master or graduate
	Formats            []CourseFormat     `json:"formats"`             // a set of preferable course formats
	Groups             []Group            `json:"groups"`              // a study groups, e.g. "BS19-04"
	Assistants         []Teacher          `json:"assistants"`          // teacher assistants of the course
	LeadingProfessor   Teacher            `json:"leading_professor"`   // e.g. who leads the lectures
	AssistantProfessor Teacher            `json:"assistant_professor"` // e.g. who leads the tutorials, might be empty

	Classes []Class `json:"classes"` // classes of the course, i.e. the course schedule
}

// Class describes a basic lesson of the Course, e.g. a couple
type Class struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Location Location      `json:"location"`
	Start    time.Time     `json:"start"`
	Duration time.Duration `json:"duration"`
	Repeats  int           `json:"repeats"`
}

// Group describes a basic students group, e.g. "BS19-04"
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Teacher describes a basic teacher with its own name and surname
type Teacher struct {
	Preferences TeacherPreferences `json:"preferences"`
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

// TimeSlot describes a particular period of time in a week
type TimeSlot struct {
	Weekday  time.Weekday      `json:"weekday"`  // a weekday of time slot
	Start    timetype.Clock    `json:"start"`    // start time of time slot
	Duration timetype.Duration `json:"duration"` // duration of a time slot
	Location Location          `json:"location"` // an optional location field, empty means "any"
}

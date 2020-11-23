package store

import (
	"time"

	"github.com/Semior001/timetype"
)

// Location describes a room or auditory where the ClassDescription is held
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
	ID              string             `json:"id"`                         // a hash derived from all others fields
	Name            string             `json:"name"`                       // the name of the course
	Program         EducationalProgram `json:"program,omitempty"`          // bachelor, master or graduate
	Formats         []CourseFormat     `json:"formats,omitempty"`          // a set of preferable course formats
	Groups          []Group            `json:"groups,omitempty"`           // a study groups, e.g. "BS19-04"
	Assistants      []Teacher          `json:"assistants,omitempty"`       // teacher assistants of the course
	PrimaryLector   Teacher            `json:"primary_lector"`             // e.g. who leads the lectures
	AssistantLector Teacher            `json:"assistant_lector,omitempty"` // e.g. who leads the tutorials, might be empty
	StudyYear       StudyYear          `json:"study_year"`

	Classes []ClassDescription `json:"classes,omitempty"` // classes of the course, i.e. the course schedule
}

// ClassDescription describes a basic lesson of the Course, e.g. a couple
type ClassDescription struct {
	ID       string        `json:"id"`
	Title    string        `json:"title"`
	Location Location      `json:"location"`
	Start    time.Time     `json:"start"`
	Duration time.Duration `json:"duration"`
	Repeats  int           `json:"repeats"`
}

// Class describes a general class with its references to groups, courses, teachers etc.
type Class struct {
	ClassDescription
	Course  Course
	Group   Group
	Teacher Teacher
}

// Group describes a basic students group, e.g. "BS19-04"
type Group struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"` // name of this group, e.g. "BS19-04"
	StudyYear StudyYear `json:"study_year"`
}

// PrepareUntrusted sets zero values to all immutable for user fields
func (g *Group) PrepareUntrusted() {
	g.ID = ""
}

// StudyYear describes a particular study year like "BS - Year 1 (Computer Science)"
type StudyYear struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TimeSlot describes a particular period of time in a week
type TimeSlot struct {
	ID       string            `json:"id"`                 // id of this time slot
	Weekday  time.Weekday      `json:"weekday"`            // a weekday of time slot
	Start    timetype.Clock    `json:"start"`              // start time of time slot
	Duration timetype.Duration `json:"duration"`           // duration of a time slot
	Location Location          `json:"location,omitempty"` // an optional location field, empty means "any"
}

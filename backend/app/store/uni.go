package store

import (
	"time"
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
	ID        string    `json:"id"`
	Name      string    `json:"name"` // name of this group, e.g. "BS19-04"
	StudyYear StudyYear `json:"study_year"`
}

// StudyYear describes a particular study year like "BS - Year 1 (Computer Science)"
type StudyYear struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

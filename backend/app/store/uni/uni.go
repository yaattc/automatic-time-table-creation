// Package uni defines interface for university organization database repository and
// provides some implementations of this interface
package uni

import "github.com/yaattc/automatic-time-table-creation/backend/app/store"

//go:generate moq -out mock_uni.go . Interface

// Interface describes database repository methods to get/set and delete groups
type Interface interface {
	AddGroup(g store.Group) (id string, err error)
	ListGroups() ([]store.Group, error)
	GetGroup(id string) (g store.Group, err error)
	DeleteGroup(id string) error

	AddStudyYear(sy store.StudyYear) (id string, err error)
	GetStudyYear(id string) (sy store.StudyYear, err error)
	DeleteStudyYear(studyYearID string) error
	ListStudyYears() ([]store.StudyYear, error)
}

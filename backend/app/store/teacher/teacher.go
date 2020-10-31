package teacher

import "github.com/yaattc/automatic-time-table-creation/backend/app/store"

//go:generate moq -out mock_teacher.go . Interface

// Interface defines repository methods for operating with teachers
type Interface interface {
	AddTeacher(teacher store.TeacherDetails) (id string, err error)
	DeleteTeacher(teacherID string) error
	ListTeachers() ([]store.TeacherDetails, error)
	GetTeacherFull(teacherID string) (store.Teacher, error)
	SetPreferences(teacherID string, pref store.TeacherPreferences) error
}

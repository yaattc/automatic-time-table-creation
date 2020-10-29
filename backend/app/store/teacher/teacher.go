package teacher

import "github.com/yaattc/automatic-time-table-creation/backend/app/store"

// Interface defines repository methods for operating with teachers
type Interface interface {
	AddTeacher(teacher store.Teacher) error
	DeleteTeacher(teacherID string) error
	ListTeachers() ([]store.Teacher, error)
	GetTeacherFull(teacherID string) (store.Teacher, error)
	SetPreferences(teacherID string, pref store.TeacherPreferences) error
}

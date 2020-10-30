package teacher

import "github.com/yaattc/automatic-time-table-creation/backend/app/store"

// Interface defines repository methods for operating with teachers
type Interface interface {
	AddTeacher(teacher store.TeacherDetails) error
	DeleteTeacher(teacherID string) error
	ListTeachers() ([]store.TeacherDetails, error)
	GetTeacherFull(teacherID string) (store.Teacher, error)
	SetPreferences(teacherID string, pref store.TeacherPreferences) error
}

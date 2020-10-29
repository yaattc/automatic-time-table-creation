package api

import (
	"net/http"

	"github.com/go-chi/chi"

	log "github.com/go-pkgz/lgr"

	"github.com/go-chi/render"
	"github.com/yaattc/automatic-time-table-creation/backend/app/rest"

	R "github.com/go-pkgz/rest"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

//nolint:unused
type private struct {
	dataService privStore
}

//go:generate moq -out mock_store.go . privStore

type privStore interface {
	AddTeacher(teacher store.Teacher) (teacherID string, err error)
	DeleteTeacher(teacherID string) error
	ListTeachers() ([]store.Teacher, error)
	GetTeacher(teacherID string) (store.Teacher, error)
	SetPreferences(teacherID string, pref store.TeacherPreferences) error
}

// POST /teacher - adds teacher
func (s *private) addTeacherCtrl(w http.ResponseWriter, r *http.Request) {
	teacher := store.Teacher{}
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, hardBodyLimit), &teacher); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't bind teacher", rest.ErrDecode)
		return
	}

	teacher.PrepareUntrusted()

	id, err := s.dataService.AddTeacher(teacher)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't add teacher", rest.ErrInternal)
		return
	}

	finalTeacher, err := s.dataService.GetTeacher(id)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't load added teacher", rest.ErrInternal)
		return
	}

	log.Printf("[DEBUG] added teacher %+v", teacher)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, finalTeacher)
}

// DELETE /teacher?id=teacherID - removes teacher
func (s *private) deleteTeacherCtrl(w http.ResponseWriter, r *http.Request) {
	teacherID := r.URL.Query().Get("id")
	if teacherID == "" {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, nil, "teacher id must be provided", rest.ErrBadRequest)
		return
	}
	if err := s.dataService.DeleteTeacher(teacherID); err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't delete teacher", rest.ErrInternal)
		return
	}

	log.Printf("[DEBUG] removed teacher %s", teacherID)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, R.JSON{"deleted": true})
}

// GET /teacher?id=teacherID - list teachers, shrink query parameter "id" to list all
func (s *private) listTeachersCtrl(w http.ResponseWriter, r *http.Request) {
	teacherID := r.URL.Query().Get("id")

	// get all teachers
	if teacherID == "" {
		ts, err := s.dataService.ListTeachers()
		if err != nil {
			rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't load teachers", rest.ErrInternal)
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, R.JSON{"teachers": ts})
		return
	}

	// get particular teacher
	teacher, err := s.dataService.GetTeacher(teacherID)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't load teacher", rest.ErrInternal)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, R.JSON{"teachers": []store.Teacher{teacher}})
}

// POST /teacher/{id}/preferences - set teacher preferences
func (s *private) setTeacherPreferencesCtrl(w http.ResponseWriter, r *http.Request) {
	teacherID := chi.URLParam(r, "id")
	pref := store.TeacherPreferences{}
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, hardBodyLimit), &pref); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't bind preferences", rest.ErrDecode)
		return
	}

	if err := s.dataService.SetPreferences(teacherID, pref); err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't set preferences for teacher", rest.ErrInternal)
		return
	}

	finalTeacher, err := s.dataService.GetTeacher(teacherID)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't load updated teacher", rest.ErrInternal)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, finalTeacher)
}

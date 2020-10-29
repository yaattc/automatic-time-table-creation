package api

import (
	"net/http"

	log "github.com/go-pkgz/lgr"

	"github.com/go-chi/render"
	"github.com/yaattc/automatic-time-table-creation/backend/app/rest"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

type private struct {
	dataService privStore
}

//go:generate moq -out mock_store.go . privStore

type privStore interface {
	AddTeacher(teacher store.Teacher) (teacherID string, err error)
	DeleteTeacher(teacherID string) error
	ListTeachers() ([]store.Teacher, error)
	GetTeacher(teacherID string) (store.Teacher, error)
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
	render.JSON(w, r, rest.JSON{"deleted": true})
}

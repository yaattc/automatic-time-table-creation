package api

import (
	"net/http"

	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/yaattc/automatic-time-table-creation/backend/app/rest"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

type uniCtrlGroup struct {
	dataService uniStore
}

//go:generate moq -out mock_uni_store.go . uniStore

type uniStore interface {
	AddGroup(name string, studyYearID string) (id string, err error)
	ListGroups() ([]store.Group, error)
	GetGroup(groupID string) (store.Group, error)
	DeleteGroup(id string) error

	AddStudyYear(name string) (id string, err error)
	GetStudyYear(id string) (sy store.StudyYear, err error)
	DeleteStudyYear(studyYearID string) error
	ListStudyYears() ([]store.StudyYear, error)

	ListTimeSlots() ([]store.TimeSlot, error)
}

// POST /group - add group
func (s *uniCtrlGroup) addGroup(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name        string `json:"name"`
		StudyYearID string `json:"study_year_id"`
	}
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, hardBodyLimit), &reqBody); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't bind group", rest.ErrDecode)
		return
	}

	id, err := s.dataService.AddGroup(reqBody.Name, reqBody.StudyYearID)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't add group", rest.ErrInternal)
		return
	}

	finalGroup, err := s.dataService.GetGroup(id)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't load added group", rest.ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, finalGroup)
}

// GET /group - list groups
func (s *uniCtrlGroup) listGroups(w http.ResponseWriter, r *http.Request) {
	tg, err := s.dataService.ListGroups()
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't load teachers", rest.ErrInternal)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, R.JSON{"groups": tg})
}

// DELETE /group?id=groupID - removes group
func (s *uniCtrlGroup) deleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("id")
	if groupID == "" {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, nil, "group id must be provided", rest.ErrBadRequest)
		return
	}
	if err := s.dataService.DeleteGroup(groupID); err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't delete teacher", rest.ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, R.JSON{"deleted": true})
}

// POST /study_year - add study year
func (s *uniCtrlGroup) addStudyYear(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name string `json:"name"`
	}
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, hardBodyLimit), &reqBody); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't bind study year", rest.ErrDecode)
		return
	}

	id, err := s.dataService.AddStudyYear(reqBody.Name)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't add study year", rest.ErrInternal)
		return
	}
	finalSy, err := s.dataService.GetStudyYear(id)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't load added study year", rest.ErrInternal)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, finalSy)
}

// GET /study_year - list study years
func (s *uniCtrlGroup) listStudyYears(w http.ResponseWriter, r *http.Request) {
	sys, err := s.dataService.ListStudyYears()
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't list study years", rest.ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, R.JSON{"study_years": sys})
}

// DELETE /study_year?id=studyYearID - remove study year
func (s *uniCtrlGroup) deleteStudyYear(w http.ResponseWriter, r *http.Request) {
	studyYearID := r.URL.Query().Get("id")
	if studyYearID == "" {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, nil, "study year's id must be provided", rest.ErrBadRequest)
		return
	}
	if err := s.dataService.DeleteStudyYear(studyYearID); err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't delete studyYear", rest.ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, R.JSON{"deleted": true})
}

// GET /time_slots - list time slots
func (s *uniCtrlGroup) listTimeSlots(w http.ResponseWriter, r *http.Request) {

}

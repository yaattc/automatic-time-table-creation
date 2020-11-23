package api

import (
	"net/http"
	"time"

	"github.com/yaattc/automatic-time-table-creation/backend/app/gen"

	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/yaattc/automatic-time-table-creation/backend/app/rest"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

//go:generate moq -out mock_sched_store.go . schedStore

type schedCtrlGroup struct {
	dataService schedStore
	genService  gen.Service
}

type schedStore interface {
	ListClasses(from time.Time, till time.Time, groupID string) ([]store.Class, error)
	ListTimeSlots() ([]store.TimeSlot, error)
}

// GET /time_slots - list time slots
func (s *schedCtrlGroup) listTimeSlots(w http.ResponseWriter, r *http.Request) {
	tsl, err := s.dataService.ListTimeSlots()
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't list time slots", rest.ErrInternal)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, R.JSON{"time_slots": tsl})
}

// POST /classes - get classes of the given group for the given period
func (s *schedCtrlGroup) listClasses(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		From    time.Time `json:"from"`
		Till    time.Time `json:"till"`
		GroupID string    `json:"group_id"`
	}
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, hardBodyLimit), &reqBody); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't bind request for list classes", rest.ErrDecode)
		return
	}

	classes, err := s.dataService.ListClasses(reqBody.From, reqBody.Till, reqBody.GroupID)
	if err != nil {
		rest.SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't load classes", rest.ErrInternal)
		return
	}

	type respBody struct {
		ID       string         `json:"id"`
		Title    string         `json:"title"`
		Location store.Location `json:"location"`
		Start    time.Time      `json:"start"`
		End      time.Time      `json:"end"`
	}

	var clDescs []respBody
	for _, cl := range classes {
		clDescs = append(clDescs, respBody{
			ID:       cl.ID,
			Title:    cl.Title,
			Location: cl.Location,
			Start:    cl.Start,
			End:      cl.Start.Add(cl.Duration),
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, R.JSON{"classes": clDescs})
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Semior001/timetype"

	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

func TestPrivate_addGroupCtrl(t *testing.T) {
	expected := store.Group{
		ID:   "00000000-0000-0000-0000-000000000003",
		Name: "B20-05",
		StudyYear: store.StudyYear{
			ID:   "00000000-0000-0000-0000-100000000003",
			Name: "BS - Year 1 (Computer Science)",
		},
	}

	var reqBody struct {
		Name        string `json:"name"`
		StudyYearID string `json:"study_year_id"`
	}
	reqBody.Name = expected.Name
	reqBody.StudyYearID = expected.StudyYear.ID

	ps := &uniStoreMock{
		AddGroupFunc: func(name string, studyYearID string) (string, error) {
			assert.Equal(t, expected.Name, name)
			assert.Equal(t, expected.StudyYear.ID, studyYearID)
			return expected.ID, nil
		},
		GetGroupFunc: func(groupID string) (store.Group, error) {
			assert.Equal(t, expected.ID, groupID)
			return expected, nil
		},
	}

	ctrl := &uniCtrlGroup{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.addGroup))
	defer ts.Close()

	b, err := json.Marshal(reqBody)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(b))
	require.NoError(t, err)

	cl := http.Client{Timeout: 5 * time.Second}
	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var actual store.Group
	err = render.DecodeJSON(resp.Body, &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestPrivate_listGroupsCtrl(t *testing.T) {
	tg := []store.Group{
		{
			ID:        "00000000-0000-0000-0000-100000000001",
			Name:      "B20-01",
			StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		},
		{
			ID:        "00000000-0000-0000-0000-100000000002",
			Name:      "B20-02",
			StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		},
		{
			ID:        "00000000-0000-0000-0000-100000000003",
			Name:      "B20-03",
			StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		},
		{
			ID:        "00000000-0000-0000-0000-100000000004",
			Name:      "B20-04",
			StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		},
		{
			ID:        "00000000-0000-0000-0000-100000000005",
			Name:      "B20-05",
			StudyYear: store.StudyYear{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		},
	}
	ps := &uniStoreMock{ListGroupsFunc: func() ([]store.Group, error) { return tg, nil }}

	ctrl := &uniCtrlGroup{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.listGroups))
	defer ts.Close()

	cl := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	type resTyp struct {
		Groups []store.Group `json:"groups"`
	}
	res := resTyp{}
	err = render.DecodeJSON(resp.Body, &res)
	require.NoError(t, err)

	assert.Equal(t, resTyp{Groups: tg}, res)
}

func TestPrivate_deleteGroupCtrl(t *testing.T) {
	id := uuid.New().String()

	ps := &uniStoreMock{
		DeleteGroupFunc: func(groupID string) error {
			assert.Equal(t, id, groupID)
			return nil
		},
	}
	ctrl := &uniCtrlGroup{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.deleteGroup))
	defer ts.Close()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s?id=%s", ts.URL, id), nil)
	require.NoError(t, err)

	cl := http.Client{Timeout: 5 * time.Second}
	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var res R.JSON
	err = render.DecodeJSON(resp.Body, &res)
	require.NoError(t, err)

	assert.Equal(t, R.JSON{"deleted": true}, res)
}

func TestPrivate_addStudyYearCtrl(t *testing.T) {
	expected := store.StudyYear{
		ID:   "00000000-0000-0000-0000-100000000003",
		Name: "BS - Year 1 (Computer Science)",
	}

	ps := &uniStoreMock{
		AddStudyYearFunc: func(name string) (string, error) {
			assert.Equal(t, expected.Name, name)
			return expected.ID, nil
		},
		GetStudyYearFunc: func(id string) (store.StudyYear, error) {
			assert.Equal(t, expected.ID, id)
			return expected, nil
		},
	}

	ctrl := &uniCtrlGroup{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.addStudyYear))
	defer ts.Close()

	var reqBody struct {
		Name string `json:"name"`
	}
	reqBody.Name = expected.Name

	b, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(b))
	require.NoError(t, err)

	cl := http.Client{Timeout: 5 * time.Second}
	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var actual store.StudyYear
	err = render.DecodeJSON(resp.Body, &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestPrivate_listStudyYearsCtrl(t *testing.T) {
	expected := []store.StudyYear{
		{ID: "00000000-0000-0000-0000-000000000001", Name: "BS - Year 1 (Computer Science)"},
		{ID: "00000000-0000-0000-0000-000000000002", Name: "MS - Year 1 (Computer Science)"},
		{ID: "00000000-0000-0000-0000-000000000003", Name: "BS - Year 2 (Computer Science)"},
		{ID: "00000000-0000-0000-0000-000000000004", Name: "MS - Year 2 (Computer Science)"},
		{ID: "00000000-0000-0000-0000-000000000005", Name: "BS - Year 3 (Computer Science)"},
	}
	ps := &uniStoreMock{ListStudyYearsFunc: func() ([]store.StudyYear, error) {
		return expected, nil
	}}

	ctrl := &uniCtrlGroup{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.listStudyYears))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	cl := http.Client{Timeout: 5 * time.Second}
	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var actual struct {
		StudyYears []store.StudyYear `json:"study_years"`
	}
	err = render.DecodeJSON(resp.Body, &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual.StudyYears)
}

func TestPrivate_deleteStudyYearCtrl(t *testing.T) {
	id := uuid.New().String()

	ps := &uniStoreMock{
		DeleteStudyYearFunc: func(studyYearID string) error {
			assert.Equal(t, id, studyYearID)
			return nil
		},
	}
	ctrl := &uniCtrlGroup{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.deleteStudyYear))
	defer ts.Close()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s?id=%s", ts.URL, id), nil)
	require.NoError(t, err)

	cl := http.Client{Timeout: 5 * time.Second}
	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var res R.JSON
	err = render.DecodeJSON(resp.Body, &res)
	require.NoError(t, err)

	assert.Equal(t, R.JSON{"deleted": true}, res)
}

func Test_uniCtrlGroup_listTimeSlots(t *testing.T) {
	expected := prepareTimeSlots()
	ps := &uniStoreMock{ListTimeSlotsFunc: func() ([]store.TimeSlot, error) {
		return expected, nil
	}}

	ctrl := &uniCtrlGroup{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.listTimeSlots))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	cl := http.Client{Timeout: 5 * time.Second}
	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var actual struct {
		TimeSlots []store.TimeSlot `json:"time_slots"`
	}
	err = render.DecodeJSON(resp.Body, &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual.TimeSlots)
}

// fixme this method is duplicated at least three times
func prepareTimeSlots() []store.TimeSlot {
	timeSlotsOnWeek := func(ts store.TimeSlot) []store.TimeSlot {
		var res []store.TimeSlot
		for i := time.Monday; i <= time.Friday; i++ {
			newTS := ts
			newTS.ID = uuid.New().String()
			newTS.Weekday = i
			res = append(res, newTS)
		}
		return res
	}
	var timeSlots []store.TimeSlot

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		Start:    timetype.NewUTCClock(9, 0, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		Start:    timetype.NewUTCClock(10, 40, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		Start:    timetype.NewUTCClock(12, 40, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		Start:    timetype.NewUTCClock(14, 20, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		Start:    timetype.NewUTCClock(16, 0, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		Start:    timetype.NewUTCClock(17, 40, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		Start:    timetype.NewUTCClock(19, 20, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	return timeSlots
}

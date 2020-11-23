package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"

	"github.com/Semior001/timetype"

	R "github.com/go-pkgz/rest"

	"github.com/go-chi/render"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

func Test_private_addTeacherCtrl(t *testing.T) {
	expected := store.Teacher{
		TeacherDetails: store.TeacherDetails{
			ID:      "this should be changed, as this field is immutable for user",
			Name:    "foo",
			Surname: "bar",
			Email:   "foo@bar.com",
			Degree:  "graduate",
			About:   "some details about teacher",
		},
		Preferences: store.TeacherPreferences{
			TimeSlots: []store.TimeSlot{
				{
					Weekday:  time.Monday,
					Start:    timetype.NewClock(20, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				},
			},
			Locations: nil,
		},
	}

	ps := &teacherStoreMock{
		AddTeacherFunc: func(teacher store.Teacher) (string, error) {
			assert.Empty(t, teacher.ID)
			assert.Empty(t, teacher.Preferences)
			expected.ID = ""
			expected.Preferences = store.TeacherPreferences{}
			assert.Equal(t, expected, teacher)
			expected.ID = uuid.New().String()
			return expected.ID, nil
		},
		GetTeacherFullFunc: func(teacherID string) (store.Teacher, error) {
			assert.Equal(t, expected.ID, teacherID)
			return expected, nil
		},
	}

	ctrl := &teacherCtrlGroup{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.addTeacherCtrl))
	defer ts.Close()

	b, err := json.Marshal(expected)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(b))
	require.NoError(t, err)

	cl := http.Client{Timeout: 5 * time.Second}
	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var actual store.Teacher
	err = render.DecodeJSON(resp.Body, &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func Test_private_deleteTeacherCtrl(t *testing.T) {
	id := uuid.New().String()

	ps := &teacherStoreMock{
		DeleteTeacherFunc: func(teacherID string) error {
			assert.Equal(t, id, teacherID)
			return nil
		},
	}
	ctrl := &teacherCtrlGroup{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.deleteTeacherCtrl))
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

func Test_private_listTeachersCtrl(t *testing.T) {
	tss := []store.Teacher{
		{
			TeacherDetails: store.TeacherDetails{

				ID:      "someFancyID",
				Name:    "foo",
				Surname: "bar",
				Email:   "foo@bar.com",
				Degree:  "graduate",
				About:   "some details about teacher",
			},
			Preferences: store.TeacherPreferences{
				TimeSlots: []store.TimeSlot{
					{
						Weekday:  time.Monday,
						Start:    timetype.NewClock(20, 0, 0, 0, time.UTC),
						Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					},
					{
						Weekday:  time.Tuesday,
						Start:    timetype.NewClock(10, 0, 0, 0, time.UTC),
						Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					},
					{
						Weekday:  time.Friday,
						Start:    timetype.NewClock(15, 0, 0, 0, time.UTC),
						Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					},
				},
				Locations: []store.Location{"108", "102", "109"},
			},
		},
		{
			TeacherDetails: store.TeacherDetails{
				ID:      "someFancyID2",
				Name:    "Nikolay",
				Surname: "Shilov",
				Email:   "s.nikolay@idontknowemail.com",
				Degree:  "graduate",
				About:   "some details about Nikolay Shilov",
			},
			Preferences: store.TeacherPreferences{
				TimeSlots: []store.TimeSlot{{
					Weekday:  time.Monday,
					Start:    timetype.NewClock(21, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				},
					{
						Weekday:  time.Tuesday,
						Start:    timetype.NewClock(13, 30, 0, 0, time.UTC),
						Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					},
					{
						Weekday:  time.Friday,
						Start:    timetype.NewClock(14, 0, 0, 0, time.UTC),
						Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					},
				},
				Locations: []store.Location{"512", "234", "312", "222"},
			},
		},
	}
	ps := &teacherStoreMock{
		ListTeachersFunc: func() (res []store.TeacherDetails, _ error) {
			for _, t := range tss {
				res = append(res, t.TeacherDetails)
			}
			return res, nil
		},
		GetTeacherFullFunc: func(teacherID string) (store.Teacher, error) {
			assert.Equal(t, tss[0].ID, teacherID)
			return tss[0], nil
		},
	}

	ctrl := &teacherCtrlGroup{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.listTeachersCtrl))
	defer ts.Close()

	cl := http.Client{Timeout: 5 * time.Second}

	// checking "get one teacher"
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?id=%s", ts.URL, tss[0].ID), nil)
	require.NoError(t, err)

	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	type resTyp struct {
		Teachers []store.TeacherDetails `json:"teachers"`
	}

	var res resTyp
	err = render.DecodeJSON(resp.Body, &res)
	require.NoError(t, err)

	assert.Equal(t, resTyp{Teachers: []store.TeacherDetails{tss[0].TeacherDetails}}, res)

	// checking "get list of teachers"

	req, err = http.NewRequest("GET", ts.URL, nil)
	require.NoError(t, err)

	resp, err = cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	res = resTyp{}
	err = render.DecodeJSON(resp.Body, &res)
	require.NoError(t, err)

	assert.Equal(t, resTyp{Teachers: []store.TeacherDetails{tss[0].TeacherDetails, tss[1].TeacherDetails}}, res)
}

func Test_private_setTeacherPreferencesCtrl(t *testing.T) {
	expected := store.Teacher{
		TeacherDetails: store.TeacherDetails{

			ID:      "someFancyID",
			Name:    "foo",
			Surname: "bar",
			Email:   "foo@bar.com",
			Degree:  "graduate",
			About:   "some details about teacher",
		},
		Preferences: store.TeacherPreferences{
			TimeSlots: []store.TimeSlot{
				{
					Weekday:  time.Monday,
					Start:    timetype.NewClock(20, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				},
				{
					Weekday:  time.Tuesday,
					Start:    timetype.NewClock(10, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				},
				{
					Weekday:  time.Friday,
					Start:    timetype.NewClock(15, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				},
			},
			Locations: []store.Location{"108", "102", "109"},
		},
	}
	ps := &teacherStoreMock{
		SetTeacherPreferencesFunc: func(teacherID string, pref store.TeacherPreferences) error {
			assert.Equal(t, expected.ID, teacherID)
			assert.Equal(t, expected.Preferences, pref)
			return nil
		},
		GetTeacherFullFunc: func(teacherID string) (store.Teacher, error) {
			assert.Equal(t, expected.ID, teacherID)
			return expected, nil
		},
	}
	ctrl := &teacherCtrlGroup{dataService: ps}
	r := chi.NewRouter()
	r.Post("/{id}", ctrl.setTeacherPreferencesCtrl)

	ts := httptest.NewServer(r)
	defer ts.Close()

	b, err := json.Marshal(expected.Preferences)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", ts.URL, expected.ID), bytes.NewReader(b))
	require.NoError(t, err)

	cl := http.Client{Timeout: 5 * time.Second}
	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var teacher store.Teacher
	err = render.DecodeJSON(resp.Body, &teacher)
	require.NoError(t, err)

	assert.Equal(t, expected, teacher)
}

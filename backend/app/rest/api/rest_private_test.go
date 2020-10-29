package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yaattc/automatic-time-table-creation/backend/app/rest"

	"github.com/go-chi/render"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

func Test_private_addTeacherCtrl(t *testing.T) {
	expected := store.Teacher{
		ID:      "this should be changed, as this field is immutable for user",
		Name:    "foo",
		Surname: "bar",
	}

	ps := &privStoreMock{
		AddTeacherFunc: func(teacher store.Teacher) (string, error) {
			assert.Empty(t, teacher.ID)
			expected.ID = ""
			assert.Equal(t, expected, teacher)
			expected.ID = uuid.New().String()
			return expected.ID, nil
		},
		GetTeacherFunc: func(teacherID string) (store.Teacher, error) {
			assert.Equal(t, expected.ID, teacherID)
			return expected, nil
		},
	}

	ctrl := &private{dataService: ps}
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

	ps := &privStoreMock{
		DeleteTeacherFunc: func(teacherID string) error {
			assert.Equal(t, id, teacherID)
			return nil
		},
	}
	ctrl := &private{dataService: ps}
	ts := httptest.NewServer(http.HandlerFunc(ctrl.deleteTeacherCtrl))
	defer ts.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s?id=%s", ts.URL, id), nil)
	require.NoError(t, err)

	cl := http.Client{Timeout: 5 * time.Second}
	resp, err := cl.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var res rest.JSON
	err = render.DecodeJSON(resp.Body, &res)
	require.NoError(t, err)

	assert.Equal(t, rest.JSON{"deleted": true}, res)
}

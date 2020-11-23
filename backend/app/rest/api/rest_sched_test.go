package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

func Test_schedCtrlGroup_listTimeSlots(t *testing.T) {
	expected := prepareTimeSlots()
	ps := &schedStoreMock{ListTimeSlotsFunc: func() ([]store.TimeSlot, error) {
		return expected, nil
	}}

	ctrl := &schedCtrlGroup{dataService: ps}
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

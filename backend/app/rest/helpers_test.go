package rest

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHourMinSec_UnmarshalJSON(t *testing.T) {
	var s struct {
		T Clock `json:"t"`
	}
	err := json.Unmarshal([]byte(`{"t": "19:24:00"}`), &s)
	require.NoError(t, err)
	assert.Equal(t, Clock(time.Date(0, time.January, 1, 19, 24, 0, 0, time.UTC)), s.T)

	s = struct {
		T Clock `json:"t"`
	}{}
	err = json.Unmarshal([]byte(`{"t": "08:00:00"}`), &s)
	require.NoError(t, err)
	assert.Equal(t, Clock(time.Date(0, time.January, 1, 8, 00, 0, 0, time.UTC)), s.T)
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	type temp struct {
		T Duration `json:"t"`
	}
	var s temp
	err := json.Unmarshal([]byte(`{"t": "1h5m3s"}`), &s)
	require.NoError(t, err)
	assert.Equal(t, Duration(time.Hour+5*time.Minute+3*time.Second), s.T)

	s = temp{}
	err = json.Unmarshal([]byte(`{"t": 3903000000000}`), &s)
	require.NoError(t, err)
	assert.Equal(t, Duration(time.Hour+5*time.Minute+3*time.Second), s.T)

	s = temp{}
	err = json.Unmarshal([]byte(`{"t": true}`), &s)
	assert.Errorf(t, err, "invalid duration")
}

func TestDuration_MarshalJSON(t *testing.T) {
	bytes, err := Duration(time.Hour + 5*time.Minute + 3*time.Second).MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte(`"1h5m3s"`), bytes)
}

func TestHourMinSec_MarshalJSON(t *testing.T) {
	bytes, err := Clock(time.Date(0, time.January, 1, 19, 24, 0, 0, time.UTC)).MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte(`"19:24:00"`), bytes)
}

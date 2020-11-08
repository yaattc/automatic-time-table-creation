package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup_PrepareUntrusted(t *testing.T) {
	unprepared := Group{
		ID:   "this should not be presented in the result",
		Name: "some group",
		StudyYear: StudyYear{
			ID:   "this is the id of the study year",
			Name: "and this is its name",
		},
	}
	prep := unprepared
	prep.PrepareUntrusted()
	unprepared.ID = ""
	assert.Equal(t, unprepared, prep)
}

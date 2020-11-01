package store

import (
	"testing"
	"time"

	"github.com/Semior001/timetype"
	"github.com/stretchr/testify/assert"
)

func TestTeacher_Empty(t *testing.T) {
	assert.True(t, TeacherPreferences{}.Empty())
	assert.False(t, TeacherPreferences{
		TimeSlots: []TimeSlot{
			{
				Weekday:  time.Monday,
				Start:    timetype.NewClock(20, 0, 0, 0, time.UTC),
				Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				Location: "room 108",
			},
			{
				Weekday:  time.Tuesday,
				Start:    timetype.NewClock(10, 0, 0, 0, time.UTC),
				Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				Location: "room 109",
			},
			{
				Weekday:  time.Friday,
				Start:    timetype.NewClock(15, 0, 0, 0, time.UTC),
				Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
				Location: "room 102",
			},
		},
		Staff:     nil,
		Locations: nil,
	}.Empty())
}

func TestTeacher_PrepareUntrusted(t *testing.T) {
	prepared := Teacher{
		TeacherDetails: TeacherDetails{
			Name:    "foo",
			Surname: "bar",
			Email:   "foo@bar.com",
			Degree:  "graduate",
			About:   "some details about teacher",
		},
	}
	toPrepare := Teacher{
		TeacherDetails: TeacherDetails{
			ID:      "00000000-0000-0000-0000-000000000001",
			Name:    "foo",
			Surname: "bar",
			Email:   "foo@bar.com",
			Degree:  "graduate",
			About:   "some details about teacher",
		},
		Preferences: TeacherPreferences{
			TimeSlots: []TimeSlot{
				{
					Weekday:  time.Monday,
					Start:    timetype.NewClock(20, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					Location: "room 108",
				},
				{
					Weekday:  time.Tuesday,
					Start:    timetype.NewClock(10, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					Location: "room 109",
				},
				{
					Weekday:  time.Friday,
					Start:    timetype.NewClock(15, 0, 0, 0, time.UTC),
					Duration: timetype.Duration(1*time.Hour + 30*time.Minute),
					Location: "room 102",
				},
			},
			Locations: []Location{"108", "102", "109"},
			Staff: []TeacherDetails{
				{
					ID:      "00000000-0000-0000-0000-000000000002",
					Name:    "Ivan",
					Surname: "Konyukhov",
					Email:   "i.konyukhov@innopolis.ru",
					Degree:  "Dr.",
					About:   "Good man",
				},
				{
					ID:      "00000000-0000-0000-0000-000000000003",
					Name:    "Nikolay",
					Surname: "Shilov",
					Email:   "n.shilov@innopolis.ru",
					Degree:  "Prof.",
					About:   "???",
				},
			},
		},
	}
	toPrepare.PrepareUntrusted()
	assert.Equal(t, prepared, toPrepare)
}

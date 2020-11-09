package gen

import (
	"fmt"
	"testing"
	"time"

	"github.com/Semior001/timetype"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

func Test_timeTable_Build(t *testing.T) {
	tt := &timeTable{}
	tt.Build(prepareTimeSlots(), prepareCourses())
	// todo asserts
}

func prepareCourses() []store.Course {
	succi := store.Teacher{
		Preferences: store.TeacherPreferences{TimeSlots: []store.TimeSlot{
			{ID: "ts9000_1", Weekday: time.Monday},
			{ID: "ts1040_1", Weekday: time.Monday},
			{ID: "ts9000_3", Weekday: time.Wednesday},
			{ID: "ts1040_3", Weekday: time.Wednesday},
			{ID: "ts1240_3", Weekday: time.Wednesday},
		}},
		TeacherDetails: store.TeacherDetails{ID: "succi"},
	}
	bobrov := store.Teacher{
		Preferences: store.TeacherPreferences{TimeSlots: []store.TimeSlot{
			{ID: "ts1040_2", Weekday: time.Tuesday},
		}},
		TeacherDetails: store.TeacherDetails{ID: "bobrov"},
	}
	kabanov := store.Teacher{
		Preferences: store.TeacherPreferences{TimeSlots: []store.TimeSlot{
			{ID: "ts1040_2", Weekday: time.Tuesday},
			{ID: "ts1240_2", Weekday: time.Tuesday},
			{ID: "ts1040_3", Weekday: time.Wednesday},
			{ID: "ts1240_3", Weekday: time.Wednesday},
		}},
		TeacherDetails: store.TeacherDetails{ID: "kabanov"},
	}
	sidorov := store.Teacher{
		Preferences: store.TeacherPreferences{TimeSlots: []store.TimeSlot{
			{ID: "ts1040_1", Weekday: time.Monday},
		}},
		TeacherDetails: store.TeacherDetails{ID: "sidorov"},
	}

	course1 := store.Course{ID: "course1", LeadingProfessor: succi, AssistantProfessor: succi}
	course2 := store.Course{ID: "course2", LeadingProfessor: bobrov, AssistantProfessor: kabanov}
	course3 := store.Course{ID: "course3", LeadingProfessor: sidorov}
	return []store.Course{course1, course2, course3}
}

func prepareTimeSlots() []store.TimeSlot {
	timeSlotsOnWeek := func(ts store.TimeSlot) []store.TimeSlot {
		var res []store.TimeSlot
		for i := time.Monday; i <= time.Friday; i++ {
			newTS := ts
			newTS.Weekday = i
			newTS.ID = fmt.Sprintf("%s_%d", ts.ID, i)
			res = append(res, newTS)
		}
		return res
	}
	var timeSlots []store.TimeSlot

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		ID:       "ts9000",
		Start:    timetype.NewUTCClock(9, 0, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		ID:       "ts1040",
		Start:    timetype.NewUTCClock(10, 40, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		ID:       "ts1240",
		Start:    timetype.NewUTCClock(12, 40, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		ID:       "ts1420",
		Start:    timetype.NewUTCClock(14, 20, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		ID:       "ts1600",
		Start:    timetype.NewUTCClock(16, 0, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		ID:       "ts1740",
		Start:    timetype.NewUTCClock(17, 40, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	timeSlots = append(timeSlots, timeSlotsOnWeek(store.TimeSlot{
		ID:       "ts1920",
		Start:    timetype.NewUTCClock(19, 20, 0, 0),
		Duration: timetype.Duration(90 * time.Minute),
	})...)

	return timeSlots
}

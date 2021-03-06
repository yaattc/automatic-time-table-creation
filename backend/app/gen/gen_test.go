package gen

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Semior001/timetype"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

func TestService_Build(t *testing.T) {
	srv := &Service{}
	crss := prepareCourses()

	res := srv.Build(BuildTimeTableRequest{
		TimeSlots: prepareTimeSlots(),
		Courses:   crss,
		From:      time.Date(2020, 11, 9, 0, 0, 0, 0, time.UTC),
		Till:      time.Date(2020, 11, 15, 0, 0, 0, 0, time.UTC),
	})
	sort.Slice(res.Classes, func(i, j int) bool {
		return strings.Compare(res.Classes[i].Title, res.Classes[j].Title) < 0
	})

	succi, bobrov, kabanov, sidorov, ivanov := prepareTeachers()
	succi.Preferences, bobrov.Preferences, kabanov.Preferences = store.TeacherPreferences{}, store.TeacherPreferences{}, store.TeacherPreferences{}
	sidorov.Preferences, ivanov.Preferences = store.TeacherPreferences{}, store.TeacherPreferences{}
	expected := []store.Class{
		{ClassDescription: store.ClassDescription{
			Title:    "course1 Lecture",
			Start:    time.Date(2020, 11, 9, 9, 0, 0, 0, time.UTC),
			Duration: 90 * time.Minute,
		}, Group: store.Group{ID: "grp1"}, Teacher: succi, Course: crss[0]},
		{ClassDescription: store.ClassDescription{
			Title:    "course1 Tutorial",
			Start:    time.Date(2020, 11, 9, 10, 40, 0, 0, time.UTC),
			Duration: 90 * time.Minute,
		}, Group: store.Group{ID: "grp1"}, Teacher: succi, Course: crss[0]},
		{ClassDescription: store.ClassDescription{
			Title:    "course2 Lecture",
			Start:    time.Date(2020, 11, 10, 10, 40, 0, 0, time.UTC),
			Duration: 90 * time.Minute,
		}, Group: store.Group{ID: "grp1"}, Teacher: bobrov, Course: crss[1]},
		{ClassDescription: store.ClassDescription{
			Title:    "course2 Tutorial",
			Start:    time.Date(2020, 11, 10, 14, 20, 0, 0, time.UTC),
			Duration: 90 * time.Minute,
		}, Group: store.Group{ID: "grp1"}, Teacher: kabanov, Course: crss[1]},
		{ClassDescription: store.ClassDescription{
			Title:    "course3 Lecture",
			Start:    time.Date(2020, 11, 9, 10, 40, 0, 0, time.UTC),
			Duration: 90 * time.Minute,
		}, Group: store.Group{ID: "grp1"}, Teacher: sidorov, Course: crss[2]},
		{ClassDescription: store.ClassDescription{
			Title:    "course4 Lecture",
			Start:    time.Date(2020, 11, 12, 10, 40, 0, 0, time.UTC),
			Duration: 90 * time.Minute,
		}, Group: store.Group{ID: "grp1"}, Teacher: kabanov, Course: crss[3]},
		{ClassDescription: store.ClassDescription{
			Title:    "course4 Tutorial",
			Start:    time.Date(2020, 11, 12, 12, 40, 0, 0, time.UTC),
			Duration: 90 * time.Minute,
		}, Group: store.Group{ID: "grp1"}, Teacher: ivanov, Course: crss[3]},
		{ClassDescription: store.ClassDescription{
			Title:    "course5 Lecture",
			Start:    time.Date(2020, 11, 11, 12, 40, 0, 0, time.UTC),
			Duration: 90 * time.Minute,
		}, Group: store.Group{ID: "grp1"}, Teacher: succi, Course: crss[4]},
		{ClassDescription: store.ClassDescription{
			Title:    "course5 Tutorial",
			Start:    time.Date(2020, 11, 11, 14, 20, 0, 0, time.UTC),
			Duration: 90 * time.Minute,
		}, Group: store.Group{ID: "grp1"}, Teacher: ivanov, Course: crss[4]},
	}
	assert.NotEmpty(t, res)

	for i := range res.Classes {
		res.Classes[i].ID = ""
		res.Classes[i].Start = res.Classes[i].Start.In(time.UTC)
		expected[i].Start = expected[i].Start.In(time.UTC)
		t.Log(res.Classes[i].ClassDescription)
		assert.Equal(t, expected[i], res.Classes[i])
	}
	t.Logf("Unused Courses: %v", res.UnusedCourses)
	assert.Empty(t, res.UnusedCourses)
}

func prepareTeachers() (store.Teacher, store.Teacher, store.Teacher, store.Teacher, store.Teacher) {
	succi := store.Teacher{
		Preferences: store.TeacherPreferences{TimeSlots: []store.TimeSlot{
			{ID: "ts0900_1", Weekday: time.Monday},
			{ID: "ts1040_1", Weekday: time.Monday},
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
			{ID: "ts0900_2", Weekday: time.Tuesday},
			{ID: "ts1040_2", Weekday: time.Tuesday},
			{ID: "ts1420_2", Weekday: time.Tuesday},

			{ID: "ts1040_3", Weekday: time.Wednesday},
			{ID: "ts1240_3", Weekday: time.Wednesday},

			{ID: "ts1040_4", Weekday: time.Thursday},
			{ID: "ts1240_4", Weekday: time.Thursday},
			{ID: "ts1420_4", Weekday: time.Thursday},
		}},
		TeacherDetails: store.TeacherDetails{ID: "kabanov"},
	}
	sidorov := store.Teacher{
		Preferences: store.TeacherPreferences{TimeSlots: []store.TimeSlot{
			{ID: "ts1040_1", Weekday: time.Monday},
		}},
		TeacherDetails: store.TeacherDetails{ID: "sidorov"},
	}
	ivanov := store.Teacher{
		Preferences: store.TeacherPreferences{TimeSlots: []store.TimeSlot{
			{ID: "ts1420_3", Weekday: time.Wednesday},
			{ID: "ts1040_4", Weekday: time.Thursday},
			{ID: "ts1240_4", Weekday: time.Thursday},
			{ID: "ts1420_4", Weekday: time.Thursday},
		}},
		TeacherDetails: store.TeacherDetails{ID: "ivanov"},
	}
	return succi, bobrov, kabanov, sidorov, ivanov
}

func prepareCourses() []store.Course {
	succi, bobrov, kabanov, sidorov, ivanov := prepareTeachers()
	course1 := store.Course{
		Groups:          []store.Group{{ID: "grp1"}},
		ID:              "course1",
		StudyYear:       store.StudyYear{ID: "sy1"},
		Name:            "course1",
		PrimaryLector:   succi,
		AssistantLector: succi,
	}
	course2 := store.Course{
		Groups:          []store.Group{{ID: "grp1"}},
		ID:              "course2",
		StudyYear:       store.StudyYear{ID: "sy1"},
		Name:            "course2",
		PrimaryLector:   bobrov,
		AssistantLector: kabanov,
	}
	course3 := store.Course{
		Groups:        []store.Group{{ID: "grp1"}},
		ID:            "course3",
		StudyYear:     store.StudyYear{ID: "sy2"},
		Name:          "course3",
		PrimaryLector: sidorov,
	}
	course4 := store.Course{
		Groups:          []store.Group{{ID: "grp1"}},
		ID:              "course4",
		StudyYear:       store.StudyYear{ID: "sy1"},
		Name:            "course4",
		PrimaryLector:   kabanov,
		AssistantLector: ivanov,
	}
	course5 := store.Course{
		Groups:          []store.Group{{ID: "grp1"}},
		ID:              "course5",
		StudyYear:       store.StudyYear{ID: "sy2"},
		Name:            "course5",
		PrimaryLector:   succi,
		AssistantLector: ivanov,
	}
	return []store.Course{course1, course2, course3, course4, course5}
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
		ID:       "ts0900",
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

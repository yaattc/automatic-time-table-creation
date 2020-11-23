package gen

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

// Service builds the timetable according the given requests
type Service struct{}

// Build the timetable according to the given data
func (s *Service) Build(req BuildTimeTableRequest) (res BuildTimeTableResult) {
	tt := timetable{}
	tt.fill(req)

	// building table
	tt.step(0)

	getCourseByIdx := func(idx int) course {
		for i, c := range tt.courses {
			if i == idx {
				return c
			}
		}
		panic("course not found")
	}

	// aggregating the results
	for dt := req.From; dt.Before(req.Till); dt = dt.AddDate(0, 0, 1) {
		for _, cell := range tt.bestResult.table[dt.Weekday()] {
			if len(cell.usedBy) < 1 {
				continue
			}

			st := time.Time(cell.slot.Start)
			stDate := time.Date(dt.Year(), dt.Month(), dt.Day(), st.Hour(), st.Minute(),
				st.Second(), st.Nanosecond(), dt.Location())
			dur := time.Duration(cell.slot.Duration)

			for _, rsrv := range cell.usedBy {
				crs := getCourseByIdx(rsrv.courseIdx)

				typ := "Tutorial"
				if rsrv.primary {
					typ = "Lecture"
				}

				res.Classes = append(res.Classes, store.Class{ClassDescription: store.ClassDescription{
					ID:       uuid.New().String(),
					Title:    fmt.Sprintf("%s %s", crs.course.Name, typ),
					Start:    stDate,
					Duration: dur,
				}})
			}
		}
	}

	var unusedCourses []string

	for _, crs := range tt.bestResult.courses {
		if !crs.status {
			unusedCourses = append(unusedCourses, crs.course.ID)
		}
	}
	res.UnusedCourses = unusedCourses

	return res
}

// BuildTimeTableResult describes the result of generating timetable
type BuildTimeTableResult struct {
	Classes       []store.Class
	UnusedCourses []string // ids of courses that were not filled successfully
}

// BuildTimeTableRequest describes parameters necessary to build
// the time table
type BuildTimeTableRequest struct {
	TimeSlots []store.TimeSlot
	Courses   []store.Course
	From      time.Time
	Till      time.Time
}

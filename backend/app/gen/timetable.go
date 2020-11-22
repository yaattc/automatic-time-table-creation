package gen

import (
	"sort"
	"time"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

type timetableCell struct {
	slot   store.TimeSlot
	usedBy *cellReservation
}

type cellReservation struct {
	courseIdx int
	primary   bool
}

type lector struct {
	teacher    store.TeacherDetails
	timeslots  []timeslotPreference
	reservedAt *lecReservation
}

type lecReservation struct {
	wd  time.Weekday
	idx int
}

type timeslotPreference struct {
	timeslotID string
	idx        int
	weekday    time.Weekday
}

type course struct {
	status          bool // is filled or not
	primaryLector   lector
	assistantLector *lector
	course          store.Course
}

type timetable struct {
	table   map[time.Weekday][]timetableCell
	courses []course

	bestResult struct {
		table   map[time.Weekday][]timetableCell
		courses []course
		score   int
	}
}

// recursive function
func (tt *timetable) step(courseIdx int) {
	if courseIdx >= len(tt.courses) {
		tt.checkAndUpgradeScore()
		return
	}
	// reserve and unreserve the slot for a given course
	reserve := func(courseIdx int, wd time.Weekday, tsIdx int, primary bool) {
		tt.table[wd][tsIdx].usedBy = &cellReservation{courseIdx: courseIdx, primary: primary}
		if primary {
			tt.courses[courseIdx].primaryLector.reservedAt = &lecReservation{wd: wd, idx: tsIdx}
			return
		}
		tt.courses[courseIdx].assistantLector.reservedAt = &lecReservation{wd: wd, idx: tsIdx}
	}
	unreserve := func(courseIdx int, wd time.Weekday, tsIdx int, primary bool) {
		tt.table[wd][tsIdx].usedBy = nil
		tt.courses[courseIdx].status = false
		if primary {
			tt.courses[courseIdx].primaryLector.reservedAt = nil
			return
		}
		tt.courses[courseIdx].assistantLector.reservedAt = nil
	}

	for _, pts := range tt.courses[courseIdx].primaryLector.timeslots {
		// if this slot is used by someone - skip it
		if tt.table[pts.weekday][pts.idx].usedBy != nil {
			continue
		}

		// try to reserve it
		reserve(courseIdx, pts.weekday, pts.idx, true)

		// if there is no assistant lector - then we filled the current course, go
		// onto the next one
		if tt.courses[courseIdx].assistantLector == nil {
			tt.courses[courseIdx].status = true
			tt.step(courseIdx + 1)
			// oops, looks like this slot cannot be used,
			// take off the reservation from it
			unreserve(courseIdx, pts.weekday, pts.idx, true)
			continue
		}

		// if we have an assistant lector - go through its timeslots
		for _, ats := range tt.courses[courseIdx].assistantLector.timeslots {
			// if this slot is not in the same day with the primary lector, or
			// it's earlier or at the same time with the primary lector's slot, or
			// it is already borrowed - skip it
			if ats.weekday != pts.weekday ||
				ats.idx <= pts.idx ||
				tt.table[ats.weekday][ats.idx].usedBy != nil {
				continue
			}

			// bingo! reserve it
			reserve(courseIdx, ats.weekday, ats.idx, false)

			// we filled the course, go onto the next one
			tt.courses[courseIdx].status = true

			// going to the next course
			tt.step(courseIdx + 1)

			// oops, looks like this slot cannot be used,
			// take off the reservation from it
			unreserve(courseIdx, ats.weekday, ats.idx, false)
		}

		unreserve(courseIdx, pts.weekday, pts.idx, true)

	}

	tt.step(courseIdx + 1)
}

func (tt *timetable) checkAndUpgradeScore() {
	score := 0
	for _, c := range tt.courses {
		if c.status {
			score++
		}
	}
	if score > tt.bestResult.score {
		tt.bestResult.table = map[time.Weekday][]timetableCell{}
		for wd, cells := range tt.table {
			tt.bestResult.table[wd] = make([]timetableCell, len(cells))
			copy(tt.bestResult.table[wd], cells)
		}
		tt.bestResult.courses = make([]course, len(tt.courses))
		copy(tt.bestResult.courses, tt.courses)
		tt.bestResult.score = score
	}
}

// sortTimeSlots the cells in the table according their start times
func (tt *timetable) sortTimeSlots() {
	for wd := range tt.table {
		sort.Slice(tt.table[wd], func(i, j int) bool {
			ti := time.Time(tt.table[wd][i].slot.Start)
			tj := time.Time(tt.table[wd][j].slot.Start)
			return ti.Before(tj)
		})
	}
}

// fill the table with data
func (tt *timetable) fill(req BuildTimeTableRequest) {
	tt.table = map[time.Weekday][]timetableCell{}

	// filling timeslots
	for _, ts := range req.TimeSlots {
		tt.table[ts.Weekday] = append(tt.table[ts.Weekday], timetableCell{slot: ts})
	}
	tt.sortTimeSlots()

	// filling courses
	for _, c := range req.Courses {
		crs := course{primaryLector: lector{teacher: c.PrimaryLector.TeacherDetails}, course: c}
		if !c.AssistantLector.Empty() {
			crs.assistantLector = &lector{teacher: c.AssistantLector.TeacherDetails}
		}

		for _, ts := range c.PrimaryLector.Preferences.TimeSlots {
			idx := 0
			for i, cell := range tt.table[ts.Weekday] {
				if cell.slot.ID == ts.ID {
					idx = i
				}
			}
			crs.primaryLector.timeslots = append(crs.primaryLector.timeslots, timeslotPreference{
				timeslotID: ts.ID,
				idx:        idx,
				weekday:    ts.Weekday,
			})
		}

		if !c.AssistantLector.Empty() {
			for _, ts := range c.AssistantLector.Preferences.TimeSlots {
				idx := 0
				for i, cell := range tt.table[ts.Weekday] {
					if cell.slot.ID == ts.ID {
						idx = i
					}
				}
				crs.assistantLector.timeslots = append(crs.assistantLector.timeslots, timeslotPreference{
					timeslotID: ts.ID,
					idx:        idx,
					weekday:    ts.Weekday,
				})
			}
		}

		tt.courses = append(tt.courses, crs)
	}
}

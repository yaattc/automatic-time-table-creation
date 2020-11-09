package gen

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

type classType string

const (
	classTypeLecture  classType = "lecture"
	classTypeTutorial classType = "tutorial"
)

type timeTableCell struct {
	timeSlot   store.TimeSlot
	pretenders []pretender

	pret *int // index to the pretender, if nil - there is no booking on this cell, if not nil - the slot is booked
}

type courseState struct {
	primaryLector struct {
		teacherID string
		used      bool
	}
	assistantLector struct {
		teacherID string
		used      bool
	}
}

type timeTable struct {
	// used for calculations
	table map[time.Weekday][]timeTableCell
	used  map[pretender]bool

	// best results
	mostSucceededResult map[time.Weekday][]timeTableCell
	score               int
	courses             map[string]courseState // map[courseID]state
}

type pretender struct {
	courseID string
	typ      classType
}

// init initializes the time table
func (tt *timeTable) init() {
	tt.table = map[time.Weekday][]timeTableCell{}
	tt.used = map[pretender]bool{}
	tt.courses = map[string]courseState{}
}

// fill the table with data
func (tt *timeTable) fill(timeSlots []store.TimeSlot, courses []store.Course) {
	appendPretender := func(wd time.Weekday, timeSlotID string, pret pretender) {
		for i := range tt.table[wd] {
			ts := tt.table[wd][i]
			if ts.timeSlot.ID == timeSlotID {
				ts.pretenders = append(ts.pretenders, pret)
				tt.table[wd][i] = ts
				return
			}
		}
	}

	// filling out initial course state
	for _, c := range courses {
		cs := courseState{}
		cs.primaryLector.teacherID = c.LeadingProfessor.ID
		cs.assistantLector.teacherID = c.AssistantProfessor.ID
		tt.courses[c.ID] = cs
	}

	// filling out time slots
	for _, ts := range timeSlots {
		tt.table[ts.Weekday] = append(tt.table[ts.Weekday], timeTableCell{timeSlot: ts})
	}

	// filling out pretenders
	for _, course := range courses {
		// filling pretenders-leading professors
		for _, ts := range course.LeadingProfessor.Preferences.TimeSlots {
			appendPretender(ts.Weekday, ts.ID, pretender{courseID: course.ID, typ: classTypeLecture})
		}
		// filling pretenders-assistant professors
		if !course.AssistantProfessor.Empty() {
			for _, ts := range course.AssistantProfessor.Preferences.TimeSlots {
				appendPretender(ts.Weekday, ts.ID, pretender{courseID: course.ID, typ: classTypeTutorial})
			}
		}
	}
}

// sortTimeSlots the cells in the table according their start times
func (tt *timeTable) sortTimeSlots() {
	for wd := range tt.table {
		sort.Slice(tt.table[wd], func(i, j int) bool {
			ti := time.Time(tt.table[wd][i].timeSlot.Start)
			tj := time.Time(tt.table[wd][j].timeSlot.Start)
			return ti.Before(tj)
		})
	}
}

// step iterates over timetable and tries to book the slot on any pretender
// backtracking algorithm
// fixme avoid recursion, too much recursion
func (tt *timeTable) step(wd time.Weekday, cellIdx int) bool {
	// if reached the and of the table, check that the timetable satisfies all conditions
	if wd >= time.Saturday {
		return tt.check()
	}
	if len(tt.table[wd]) <= cellIdx {
		return tt.step(wd+1, 0)
	}

	// looking for pretenders on this time slot
	for prtIdx, prt := range tt.table[wd][cellIdx].pretenders {
		// if this pretender is already used in another time slot - skip
		// todo also check that **professor** is not used at this time slot in the other subject or smth.
		if tt.used[prt] {
			continue
		}

		// try to book this time slot on the pretender
		tt.table[wd][cellIdx].pret = &prtIdx
		tt.used[prt] = true

		if tt.step(wd, cellIdx+1) {
			return true
		}

		// oops, book the pretender to this time slot was not succeed, rollback
		tt.table[wd][cellIdx].pret = nil
		tt.used[prt] = false
	}
	return tt.step(wd, cellIdx+1)
}

// check the timetable - does it satisfy all conditions
func (tt *timeTable) check() bool {
	currCoursesState := tt.getEmptyCoursesState()
	tt.fillCoursesState(currCoursesState, tt.table)

	// calculating overall state score
	score := 0
	for _, courseState := range currCoursesState {
		if courseState.primaryLector.used && courseState.assistantLector.used {
			score++
		}
	}

	// updating best result
	if score > tt.score {
		tt.score = score
		bestTable := map[time.Weekday][]timeTableCell{}
		// updating the best table
		for wd, cells := range tt.table {
			bestTable[wd] = []timeTableCell{}
			for _, cell := range cells {
				newCell := timeTableCell{timeSlot: cell.timeSlot}
				if cell.pret != nil {
					newPret := *cell.pret
					newCell.pret = &newPret
				}
				if len(cell.pretenders) > 0 {
					newCell.pretenders = make([]pretender, len(cell.pretenders))
					copy(newCell.pretenders, cell.pretenders)
				}

				bestTable[wd] = append(bestTable[wd], newCell)
			}
		}
		tt.mostSucceededResult = bestTable
		tt.courses = currCoursesState
	}

	// the best score (when all conditions are satisfied) is when all courses are filled in the table
	// todo also check the order of subjects in a weekdays
	return len(tt.courses) == score
}

type BuildRequest struct {
	TimeSlots []store.TimeSlot
	Courses   []store.Course
	From      time.Time
	Till      time.Time
}

type BuildResult struct {
	Classes       []store.Class
	UnusedCourses []string // ids of courses that couldn't be used in the time table
}

// Build the schedule
func (tt *timeTable) Build(req BuildRequest) BuildResult {
	tt.init()
	tt.fill(req.TimeSlots, req.Courses)
	tt.sortTimeSlots()

	// starting generation
	tt.step(time.Monday, 0)

	getCourseByID := func(courseID string) store.Course {
		for _, course := range req.Courses {
			if course.ID == courseID {
				return course
			}
		}
		return store.Course{}
	}

	// if we got any result
	if tt.mostSucceededResult != nil {
		var res BuildResult
		for dt := req.From; dt.Before(req.Till); dt = dt.AddDate(0, 0, 1) {
			for _, cell := range tt.mostSucceededResult[dt.Weekday()] {
				if cell.pret == nil {
					continue
				}
				st := time.Time(cell.timeSlot.Start)
				stDate := time.Date(dt.Year(), dt.Month(), dt.Day(), st.Hour(), st.Minute(),
					st.Second(), st.Nanosecond(), dt.Location())
				dur := time.Duration(cell.timeSlot.Duration)

				pret := cell.pretenders[*cell.pret]

				res.Classes = append(res.Classes, store.Class{
					ID:       uuid.New().String(),
					Title:    fmt.Sprintf("%s %s", getCourseByID(pret.courseID).Name, pret.typ),
					Start:    stDate,
					Duration: dur,
				})
			}
		}
		coursesState := tt.getEmptyCoursesState()
		tt.fillCoursesState(coursesState, tt.mostSucceededResult)
		// calculating overall state score
		for courseID, courseState := range coursesState {
			if !courseState.primaryLector.used || !courseState.assistantLector.used {
				res.UnusedCourses = append(res.UnusedCourses, courseID)
			}
		}
		return res
	}
	return BuildResult{}
}

func (tt *timeTable) getEmptyCoursesState() map[string]courseState {
	currCoursesState := map[string]courseState{}

	// copying states
	for courseID, state := range tt.courses {
		// filling state to zero
		state.primaryLector.used = false
		state.assistantLector.used = false
		// copying state
		currCoursesState[courseID] = state
	}

	return currCoursesState
}

func (tt *timeTable) fillCoursesState(st map[string]courseState, tbl map[time.Weekday][]timeTableCell) {
	// checking and updating course states
	for _, cells := range tbl {
		for _, cell := range cells {
			// if this slot booked, mark the class as used
			if cell.pret != nil {
				prt := cell.pretenders[*cell.pret]

				// updating the state
				courseState := st[prt.courseID]
				switch prt.typ {
				case classTypeLecture:
					courseState.primaryLector.used = true
				case classTypeTutorial:
					courseState.assistantLector.used = true
				}
				st[prt.courseID] = courseState
			}
		}
	}
}

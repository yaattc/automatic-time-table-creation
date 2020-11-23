package sched

import (
	"time"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

// Interface describes methods for manipulating the timetable
type Interface interface {
	ListClasses(from time.Time, till time.Time, groupID string) ([]store.Class, error)
}

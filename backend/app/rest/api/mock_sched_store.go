// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package api

import (
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
	"sync"
	"time"
)

// Ensure, that schedStoreMock does implement schedStore.
// If this is not the case, regenerate this file with moq.
var _ schedStore = &schedStoreMock{}

// schedStoreMock is a mock implementation of schedStore.
//
//     func TestSomethingThatUsesschedStore(t *testing.T) {
//
//         // make and configure a mocked schedStore
//         mockedschedStore := &schedStoreMock{
//             AddClassesFunc: func(classes []store.Class) error {
// 	               panic("mock out the AddClasses method")
//             },
//             GetCourseFunc: func(id string) (store.Course, error) {
// 	               panic("mock out the GetCourse method")
//             },
//             ListClassesFunc: func(from time.Time, till time.Time, groupID string) ([]store.Class, error) {
// 	               panic("mock out the ListClasses method")
//             },
//             ListCoursesFunc: func() ([]store.Course, error) {
// 	               panic("mock out the ListCourses method")
//             },
//             ListTimeSlotsFunc: func() ([]store.TimeSlot, error) {
// 	               panic("mock out the ListTimeSlots method")
//             },
//         }
//
//         // use mockedschedStore in code that requires schedStore
//         // and then make assertions.
//
//     }
type schedStoreMock struct {
	// AddClassesFunc mocks the AddClasses method.
	AddClassesFunc func(classes []store.Class) error

	// GetCourseFunc mocks the GetCourse method.
	GetCourseFunc func(id string) (store.Course, error)

	// ListClassesFunc mocks the ListClasses method.
	ListClassesFunc func(from time.Time, till time.Time, groupID string) ([]store.Class, error)

	// ListCoursesFunc mocks the ListCourses method.
	ListCoursesFunc func() ([]store.Course, error)

	// ListTimeSlotsFunc mocks the ListTimeSlots method.
	ListTimeSlotsFunc func() ([]store.TimeSlot, error)

	// calls tracks calls to the methods.
	calls struct {
		// AddClasses holds details about calls to the AddClasses method.
		AddClasses []struct {
			// Classes is the classes argument value.
			Classes []store.Class
		}
		// GetCourse holds details about calls to the GetCourse method.
		GetCourse []struct {
			// ID is the id argument value.
			ID string
		}
		// ListClasses holds details about calls to the ListClasses method.
		ListClasses []struct {
			// From is the from argument value.
			From time.Time
			// Till is the till argument value.
			Till time.Time
			// GroupID is the groupID argument value.
			GroupID string
		}
		// ListCourses holds details about calls to the ListCourses method.
		ListCourses []struct {
		}
		// ListTimeSlots holds details about calls to the ListTimeSlots method.
		ListTimeSlots []struct {
		}
	}
	lockAddClasses    sync.RWMutex
	lockGetCourse     sync.RWMutex
	lockListClasses   sync.RWMutex
	lockListCourses   sync.RWMutex
	lockListTimeSlots sync.RWMutex
}

// AddClasses calls AddClassesFunc.
func (mock *schedStoreMock) AddClasses(classes []store.Class) error {
	if mock.AddClassesFunc == nil {
		panic("schedStoreMock.AddClassesFunc: method is nil but schedStore.AddClasses was just called")
	}
	callInfo := struct {
		Classes []store.Class
	}{
		Classes: classes,
	}
	mock.lockAddClasses.Lock()
	mock.calls.AddClasses = append(mock.calls.AddClasses, callInfo)
	mock.lockAddClasses.Unlock()
	return mock.AddClassesFunc(classes)
}

// AddClassesCalls gets all the calls that were made to AddClasses.
// Check the length with:
//     len(mockedschedStore.AddClassesCalls())
func (mock *schedStoreMock) AddClassesCalls() []struct {
	Classes []store.Class
} {
	var calls []struct {
		Classes []store.Class
	}
	mock.lockAddClasses.RLock()
	calls = mock.calls.AddClasses
	mock.lockAddClasses.RUnlock()
	return calls
}

// GetCourse calls GetCourseFunc.
func (mock *schedStoreMock) GetCourse(id string) (store.Course, error) {
	if mock.GetCourseFunc == nil {
		panic("schedStoreMock.GetCourseFunc: method is nil but schedStore.GetCourse was just called")
	}
	callInfo := struct {
		ID string
	}{
		ID: id,
	}
	mock.lockGetCourse.Lock()
	mock.calls.GetCourse = append(mock.calls.GetCourse, callInfo)
	mock.lockGetCourse.Unlock()
	return mock.GetCourseFunc(id)
}

// GetCourseCalls gets all the calls that were made to GetCourse.
// Check the length with:
//     len(mockedschedStore.GetCourseCalls())
func (mock *schedStoreMock) GetCourseCalls() []struct {
	ID string
} {
	var calls []struct {
		ID string
	}
	mock.lockGetCourse.RLock()
	calls = mock.calls.GetCourse
	mock.lockGetCourse.RUnlock()
	return calls
}

// ListClasses calls ListClassesFunc.
func (mock *schedStoreMock) ListClasses(from time.Time, till time.Time, groupID string) ([]store.Class, error) {
	if mock.ListClassesFunc == nil {
		panic("schedStoreMock.ListClassesFunc: method is nil but schedStore.ListClasses was just called")
	}
	callInfo := struct {
		From    time.Time
		Till    time.Time
		GroupID string
	}{
		From:    from,
		Till:    till,
		GroupID: groupID,
	}
	mock.lockListClasses.Lock()
	mock.calls.ListClasses = append(mock.calls.ListClasses, callInfo)
	mock.lockListClasses.Unlock()
	return mock.ListClassesFunc(from, till, groupID)
}

// ListClassesCalls gets all the calls that were made to ListClasses.
// Check the length with:
//     len(mockedschedStore.ListClassesCalls())
func (mock *schedStoreMock) ListClassesCalls() []struct {
	From    time.Time
	Till    time.Time
	GroupID string
} {
	var calls []struct {
		From    time.Time
		Till    time.Time
		GroupID string
	}
	mock.lockListClasses.RLock()
	calls = mock.calls.ListClasses
	mock.lockListClasses.RUnlock()
	return calls
}

// ListCourses calls ListCoursesFunc.
func (mock *schedStoreMock) ListCourses() ([]store.Course, error) {
	if mock.ListCoursesFunc == nil {
		panic("schedStoreMock.ListCoursesFunc: method is nil but schedStore.ListCourses was just called")
	}
	callInfo := struct {
	}{}
	mock.lockListCourses.Lock()
	mock.calls.ListCourses = append(mock.calls.ListCourses, callInfo)
	mock.lockListCourses.Unlock()
	return mock.ListCoursesFunc()
}

// ListCoursesCalls gets all the calls that were made to ListCourses.
// Check the length with:
//     len(mockedschedStore.ListCoursesCalls())
func (mock *schedStoreMock) ListCoursesCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockListCourses.RLock()
	calls = mock.calls.ListCourses
	mock.lockListCourses.RUnlock()
	return calls
}

// ListTimeSlots calls ListTimeSlotsFunc.
func (mock *schedStoreMock) ListTimeSlots() ([]store.TimeSlot, error) {
	if mock.ListTimeSlotsFunc == nil {
		panic("schedStoreMock.ListTimeSlotsFunc: method is nil but schedStore.ListTimeSlots was just called")
	}
	callInfo := struct {
	}{}
	mock.lockListTimeSlots.Lock()
	mock.calls.ListTimeSlots = append(mock.calls.ListTimeSlots, callInfo)
	mock.lockListTimeSlots.Unlock()
	return mock.ListTimeSlotsFunc()
}

// ListTimeSlotsCalls gets all the calls that were made to ListTimeSlots.
// Check the length with:
//     len(mockedschedStore.ListTimeSlotsCalls())
func (mock *schedStoreMock) ListTimeSlotsCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockListTimeSlots.RLock()
	calls = mock.calls.ListTimeSlots
	mock.lockListTimeSlots.RUnlock()
	return calls
}

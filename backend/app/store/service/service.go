// Package service wraps user interfaces with common logic unrelated to any particular user implementation.
// All consumers should be using service.DataStore and not the naked repositories!
package service

import (
	"crypto/sha1" // nolint
	"log"
	"time"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/sched"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/uni"

	"github.com/google/uuid"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/teacher"

	"github.com/go-pkgz/auth/token"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
	"golang.org/x/crypto/bcrypt"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/user"

	"github.com/pkg/errors"
)

// DataStore wraps all stores with common and additional methods
// todo looks ugly, rewrite
type DataStore struct {
	UserRepository    user.Interface
	TeacherRepository teacher.Interface
	UniOrgRepository  uni.Interface
	SchedRepository   sched.Interface
	BCryptCost        int
}

// AddTeacher to the database
func (s *DataStore) AddTeacher(teacher store.Teacher) (teacherID string, err error) {
	if teacher.ID == "" {
		teacher.ID = uuid.New().String()
	}
	if teacherID, err = s.TeacherRepository.AddTeacher(teacher.TeacherDetails); err != nil {
		return "", errors.Wrapf(err, "failed to add teacher %s %s to database", teacher.Name, teacher.Surname)
	}
	if !teacher.Preferences.Empty() {
		if err = s.TeacherRepository.SetPreferences(teacher.ID, teacher.Preferences); err != nil {
			return "", errors.Wrapf(err,
				"failed to set preferences for teacher %s %s during the addition", teacher.Name, teacher.Surname)
		}
	}
	return teacherID, nil
}

// DeleteTeacher from the database by its id
func (s *DataStore) DeleteTeacher(teacherID string) error {
	return errors.Wrapf(s.TeacherRepository.DeleteTeacher(teacherID), "failed to delete teacher %s", teacherID)
}

// ListTeachers returns all teachers that are registered in the database
func (s *DataStore) ListTeachers() ([]store.TeacherDetails, error) {
	return s.TeacherRepository.ListTeachers()
}

// GetTeacherFull returns all data about the requested teacher, including teacher preferences
func (s *DataStore) GetTeacherFull(teacherID string) (store.Teacher, error) {
	return s.TeacherRepository.GetTeacherFull(teacherID)
}

// SetTeacherPreferences sets preferences for the given teacher
func (s *DataStore) SetTeacherPreferences(teacherID string, pref store.TeacherPreferences) error {
	return s.TeacherRepository.SetPreferences(teacherID, pref)
}

// GetUserEmail returns the email of the specified user
func (s *DataStore) GetUserEmail(id string) (email string, err error) {
	u, err := s.UserRepository.GetUser(id)
	//goland:noinspection GoNilness
	return u.Email, errors.Wrapf(err, "failed to read email of %s", id)
}

// GetUserPrivs returns the list of privileges of the specified user
func (s *DataStore) GetUserPrivs(id string) (privs []store.Privilege, err error) {
	u, err := s.UserRepository.GetUser(id)
	//goland:noinspection GoNilness
	return u.Privileges, errors.Wrapf(err, "failed to read privs of %s", id)
}

// CheckUserCredentials with the given username and password
func (s *DataStore) CheckUserCredentials(email string, password string) (ok bool, err error) {
	userpwd, err := s.UserRepository.GetPasswordHash(email)
	if err != nil {
		return false, errors.Wrapf(err, "failed to validate user")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(userpwd), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Printf("[DEBUG] wrong password for %s", email)
			return false, nil
		}
		return false, err
	}
	return true, err
}

// AddUser to the database, hash its password and give it an ID, if needed
func (s *DataStore) AddUser(user store.User, password string) (id string, err error) {
	// hashing password
	b, err := bcrypt.GenerateFromPassword([]byte(password), s.BCryptCost)
	if err != nil {
		return "", errors.Wrapf(err, "failed to hash %s user's password with bcrypt", user.Email)
	}
	// adding id
	if user.ID == "" {
		user.ID = "local_" + token.HashID(sha1.New(), user.Email) // nolint
	}

	id, err = s.UserRepository.AddUser(user, string(b), false)
	return id, errors.Wrapf(err, "failed to add user %s to database", user.ID)
}

// RegisterAdmin in the database
func (s *DataStore) RegisterAdmin(email string, password string) (id string, err error) {
	// hashing password
	b, err := bcrypt.GenerateFromPassword([]byte(password), s.BCryptCost)
	if err != nil {
		return "", errors.Wrapf(err, "failed to hash %s user's password with bcrypt", email)
	}
	u := store.User{
		ID:         "local_" + token.HashID(sha1.New(), email), // nolint
		Email:      email,
		Privileges: []store.Privilege{store.PrivReadUsers, store.PrivEditUsers, store.PrivListUsers, store.PrivAddUsers},
	}
	log.Printf("[INFO] trying to register admin with %+v and pwd %s", u, password)
	if id, err = s.UserRepository.AddUser(u, string(b), true); err != nil {
		return "", errors.Wrapf(err, "failed to add user %s to database", u.ID)
	}
	return id, nil
}

// AddGroup to the database
func (s *DataStore) AddGroup(name string, studyYearID string) (id string, err error) {
	g := store.Group{ID: uuid.New().String(), Name: name, StudyYear: store.StudyYear{ID: studyYearID}}
	id, err = s.UniOrgRepository.AddGroup(g)
	return id, errors.Wrapf(err, "failed to add group with name %s", name)
}

// GetGroup from the database
func (s *DataStore) GetGroup(id string) (store.Group, error) {
	g, err := s.UniOrgRepository.GetGroup(id)
	return g, errors.Wrapf(err, "failed to get group with id %s", id)
}

// ListGroups registered in the database
func (s *DataStore) ListGroups() ([]store.Group, error) {
	g, err := s.UniOrgRepository.ListGroups()
	return g, errors.Wrap(err, "failed to list groups")
}

// DeleteGroup from the database
func (s *DataStore) DeleteGroup(id string) error {
	return errors.Wrapf(s.UniOrgRepository.DeleteGroup(id), "failed to delete group %s", id)
}

// AddStudyYear to the database
func (s *DataStore) AddStudyYear(name string) (id string, err error) {
	sy := store.StudyYear{Name: name}
	if sy.ID == "" {
		sy.ID = uuid.New().String()
	}
	return s.UniOrgRepository.AddStudyYear(sy)
}

// GetStudyYear by its id
func (s *DataStore) GetStudyYear(id string) (sy store.StudyYear, err error) {
	return s.UniOrgRepository.GetStudyYear(id)
}

// DeleteStudyYear by its id
func (s *DataStore) DeleteStudyYear(studyYearID string) error {
	return s.UniOrgRepository.DeleteStudyYear(studyYearID)
}

// ListStudyYears that are registered in the database
func (s *DataStore) ListStudyYears() ([]store.StudyYear, error) {
	return s.UniOrgRepository.ListStudyYears()
}

// AddCourse to the database
func (s *DataStore) AddCourse(course store.Course) (id string, err error) {
	if course.ID == "" {
		course.ID = uuid.New().String()
	}
	return s.UniOrgRepository.AddCourse(course)
}

// GetCourse by id
func (s *DataStore) GetCourse(id string) (store.Course, error) {
	crs, err := s.UniOrgRepository.GetCourseDetails(id)
	if err != nil {
		return store.Course{}, errors.Wrapf(err, "failed to get details for course %s", id)
	}
	// loading teachers
	if crs.PrimaryLector, err = s.TeacherRepository.GetTeacherFull(crs.PrimaryLector.ID); err != nil {
		return store.Course{}, errors.Wrapf(err, "failed to load primary lector %s for course %s",
			crs.PrimaryLector.ID, id)
	}

	// if assistant lector is assumed for this course
	if crs.AssistantLector.ID != "" {
		if crs.AssistantLector, err = s.TeacherRepository.GetTeacherFull(crs.AssistantLector.ID); err != nil {
			return store.Course{}, errors.Wrapf(err, "failed to load assistant lector %s for course %s",
				crs.AssistantLector.ID, id)
		}
	}

	for taIdx, ta := range crs.Assistants {
		if crs.Assistants[taIdx], err = s.TeacherRepository.GetTeacherFull(ta.ID); err != nil {
			return store.Course{}, errors.Wrapf(err, "failed to load TA %s for course %s",
				ta.ID, id)
		}
	}

	return crs, nil
}

// ListTimeSlots that are registered in the database
func (s *DataStore) ListTimeSlots() ([]store.TimeSlot, error) {
	return s.UniOrgRepository.ListTimeSlots()
}

// ListClasses in the given period for the given group
func (s *DataStore) ListClasses(from time.Time, till time.Time, groupID string) ([]store.Class, error) {
	cls, err := s.SchedRepository.ListClasses(from, till, groupID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list classes details from %v to %v for group %s")
	}
	for clIdx := range cls {
		cl := &cls[clIdx]

		if cl.Course, err = s.UniOrgRepository.GetCourseDetails(cl.Course.ID); err != nil {
			return nil, errors.Wrapf(err, "failed to get course details for course %s", cl.Group.ID)
		}
		if cl.Group, err = s.UniOrgRepository.GetGroup(cl.Group.ID); err != nil {
			return nil, errors.Wrapf(err, "failed to get group details for group %s", cl.Group.ID)
		}
		if cl.Teacher.TeacherDetails, err = s.TeacherRepository.GetTeacherDetails(cl.Teacher.ID); err != nil {
			return nil, errors.Wrapf(err, "failed to get teacher details for %s", cl.Teacher.ID)
		}
	}
	return cls, nil
}

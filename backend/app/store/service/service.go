// Package service wraps user interfaces with common logic unrelated to any particular user implementation.
// All consumers should be using service.DataStore and not the naked repositories!
package service

import (
	"crypto/sha1" // nolint
	"log"

	"github.com/google/uuid"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/teacher"

	"github.com/go-pkgz/auth/token"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
	"golang.org/x/crypto/bcrypt"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/user"

	"github.com/pkg/errors"
)

// DataStore wraps all stores with common and additional methods
type DataStore struct {
	UserRepository    user.Interface
	TeacherRepository teacher.Interface
	BCryptCost        int
}

// AddTeacher to the database
func (s *DataStore) AddTeacher(teacher store.Teacher) (teacherID string, err error) {
	if teacher.ID == "" {
		teacher.ID = uuid.New().String()
	}
	if err := s.TeacherRepository.AddTeacher(teacher); err != nil {
		return "", errors.Wrapf(err, "failed to add teacher %s %s to database", teacher.Name, teacher.Surname)
	}
	return teacher.ID, nil
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
	if err != nil {
		return "", errors.Wrapf(err, "failed to read email of %s", id)
	}
	return u.Email, nil
}

// GetUserPrivs returns the list of privileges of the specified user
func (s *DataStore) GetUserPrivs(id string) (privs []store.Privilege, err error) {
	u, err := s.UserRepository.GetUser(id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read privs of %s", id)
	}
	return u.Privileges, nil
}

// CheckUserCredentials with the given username and password
func (s *DataStore) CheckUserCredentials(email string, password string) (ok bool, err error) {
	userpwd, err := s.UserRepository.GetPasswordHash(email)
	if err != nil {
		return false, errors.Wrapf(err, "failed to validate user")
	}
	err = bcrypt.CompareHashAndPassword([]byte(userpwd), []byte(password))
	return err == nil, err
}

// AddUser to the database, hash its password and give it an ID, if needed
func (s *DataStore) AddUser(user store.User, password string) (err error) {
	// hashing password
	b, err := bcrypt.GenerateFromPassword([]byte(password), s.BCryptCost)
	if err != nil {
		return errors.Wrapf(err, "failed to hash %s user's password with bcrypt", user.Email)
	}
	// adding id
	if user.ID == "" {
		user.ID = "local_" + token.HashID(sha1.New(), user.Email) // nolint // fixme
	}
	return errors.Wrapf(s.UserRepository.AddUser(user, string(b), false), "failed to add user %s to database", user.ID)
}

// RegisterAdmin in the database
func (s *DataStore) RegisterAdmin(email string, password string) error {
	// hashing password
	b, err := bcrypt.GenerateFromPassword([]byte(password), s.BCryptCost)
	if err != nil {
		return errors.Wrapf(err, "failed to hash %s user's password with bcrypt", email)
	}
	u := store.User{
		ID:         "local_" + token.HashID(sha1.New(), email), // nolint // fixme
		Email:      email,
		Privileges: []store.Privilege{store.PrivReadUsers, store.PrivEditUsers, store.PrivListUsers, store.PrivAddUsers},
	}
	log.Printf("[INFO] trying to register admin with %+v and pwd %s", u, password)
	return errors.Wrapf(s.UserRepository.AddUser(u, string(b), true), "failed to add user %s to database", u.ID)
}

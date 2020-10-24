// Package service wraps engine interfaces with common logic unrelated to any particular engine implementation.
// All consumers should be using service.DataStore and not the naked engine!
package service

import (
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
	"golang.org/x/crypto/bcrypt"

	"github.com/yaattc/automatic-time-table-creation/backend/app/store/engine"

	"github.com/pkg/errors"
)

// DataStore wraps all stores with common and additional methods
type DataStore struct {
	Engine engine.Interface
}

// GetUserEmail returns the email of the specified user
func (s *DataStore) GetUserEmail(id string) (email string, err error) {
	u, err := s.Engine.GetUser(id)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read email of %s", id)
	}
	return u.Email, nil
}

// GetUserPrivs returns the list of privileges of the specified user
func (s *DataStore) GetUserPrivs(id string) (privs []store.Privilege, err error) {
	u, err := s.Engine.GetUser(id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read privs of %s", id)
	}
	return u.Privileges, nil
}

// CheckUserCredentials with the given username and pwd
func (s *DataStore) CheckUserCredentials(email string, pwd string) (ok bool, err error) {
	userpwd, err := s.Engine.GetPasswordHash(email)
	if err != nil {
		return false, errors.Wrapf(err, "failed to validate user")
	}
	err = bcrypt.CompareHashAndPassword([]byte(userpwd), []byte(pwd))
	return err == nil, err
}

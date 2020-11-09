// Package user provides implementations for Interface for the database user repository.
// All consumers should be using service.DataStore and not the naked repositories!
package user

import (
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

//go:generate moq -out mock_user.go . Interface

// Interface defines methods to repository, and fetch models
type Interface interface {
	GetUser(id string) (u store.User, err error)
	GetPasswordHash(email string) (pwd string, err error)
	AddUser(user store.User, pwd string, ignoreIfExists bool) (id string, err error)
}

package engine

import "github.com/yaattc/automatic-time-table-creation/backend/app/store"

// Interface defines methods to store, and fetch models
type Interface interface {
	GetUser(id string) (u store.User, err error)
	GetPasswordHash(email string) (pwd string, err error)
}

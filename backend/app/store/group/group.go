package group

import "github.com/yaattc/automatic-time-table-creation/backend/app/store"

//go:generate moq -out mock_group.go . Interface

// Interface describes database repository methods to get/set and delete groups
type Interface interface {
	AddGroup(groupID string, name string) (id string, err error)
	ListGroups() ([]store.Group, error)
	DeleteGroup(id string) error
}

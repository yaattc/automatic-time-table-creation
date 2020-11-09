package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_HasPrivilege(t *testing.T) {
	assert.True(t, User{Privileges: []Privilege{PrivAddUsers, PrivListUsers}}.HasPrivilege(string(PrivListUsers)))
	assert.False(t, User{Privileges: []Privilege{PrivAddUsers, PrivListUsers}}.HasPrivilege(string(PrivReadUsers)))
}

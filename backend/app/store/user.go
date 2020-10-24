// Package store defines models and global methods to operate these models
package store

// User describes a basic user
type User struct {
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	Admin      bool     `json:"admin"`
	Privileges []string `json:"privileges"`
}

// HasPrivilege checks whether user has the defined privilege
func (u User) HasPrivilege(priv string) bool {
	for _, p := range u.Privileges {
		if p == priv {
			return true
		}
	}
	return false
}

// Package store defines models and global methods to operate these models
package store

// Privilege describes a user's privilege to perform some action
type Privilege string

// Default user privileges
const (
	PrivAddUsers  Privilege = "add_users"
	PrivReadUsers Privilege = "read_users"
	PrivEditUsers Privilege = "edit_users"
	PrivListUsers Privilege = "list_users"
)

// User describes a basic user
type User struct {
	ID         string      `json:"id"`
	Email      string      `json:"email"`
	Privileges []Privilege `json:"privileges"`
}

// HasPrivilege checks whether user has the defined privilege
func (u User) HasPrivilege(priv string) bool {
	for _, p := range u.Privileges {
		if string(p) == priv {
			return true
		}
	}
	return false
}

// PrivsToStr converts privileges slice to string slice to simplify embedding
func PrivsToStr(privs []Privilege) []string {
	var res []string
	for _, p := range privs {
		res = append(res, string(p))
	}
	return res
}

// StrToPrivs converts string slice to privileges to simplify embedding
func StrToPrivs(s []string) []Privilege {
	var res []Privilege
	// converting string-privileges to enums
	for _, p := range s {
		res = append(res, Privilege(p))
	}
	return res
}

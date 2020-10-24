package rest

import (
	"net/http"

	"github.com/go-pkgz/auth/token"
	"github.com/pkg/errors"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
)

// MustGetUserInfo fails if can't extract user data from the request.
// should be called from authed controllers only
func MustGetUserInfo(r *http.Request) store.User {
	user, err := GetUserInfo(r)
	if err != nil {
		panic(err)
	}
	return user
}

// GetUserInfo returns user from request context
func GetUserInfo(r *http.Request) (user store.User, err error) {
	u, err := token.GetUserInfo(r)
	if err != nil {
		return store.User{}, errors.Wrap(err, "can't extract user info from the token")
	}

	return store.User{
		ID:         u.ID,
		Email:      u.Email,
		Admin:      u.IsAdmin(),
		Privileges: u.SliceAttr("privileges"),
	}, nil
}

// GetUserOrEmpty attempts to get user info from request and returns empty object if failed
func GetUserOrEmpty(r *http.Request) store.User {
	user, err := GetUserInfo(r)
	if err != nil {
		return store.User{}
	}
	return user
}

// SetUserInfo sets user into request context
func SetUserInfo(r *http.Request, user store.User) *http.Request {
	u := token.User{
		ID:       user.ID,
		Email:    user.Email,
		Audience: "attc",
	}
	u.SetSliceAttr("privileges", user.Privileges)
	u.SetAdmin(user.Admin)
	return token.SetUserInfo(r, u)
}

package cmd

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store/service"
	"github.com/yaattc/automatic-time-table-creation/backend/app/store/user"
)

// AddUser adds user to the database with the specified user data
type AddUser struct {
	Location string `long:"location" env:"LOCATION" description:"location of all db files" required:"true"`
	User     struct {
		ID         string `long:"id" env:"ID" required:"false"`
		Email      string `long:"email" env:"EMAIL" required:"true"`
		Password   string `long:"password" env:"PASSWORD" required:"true" `
		Privileges string `long:"privileges" env:"PRIVILEGES" required:"true" description:"JSON-specified list of privileges"`
	} `group:"user" namespace:"user" env-namespace:"USER"`

	BCryptCost int    `long:"bcrypt_cost" env:"BCRYPT_COST" description:"bcrypt cost for hashing user password" default:"10"`
	DBConnStr  string `long:"db_conn_str" env:"DB_CONN_STR" required:"true" description:"connection string to db"`

	CommonOpts
}

// Execute runs http web server
func (a *AddUser) Execute(_ []string) error {
	pgpool, pgconf, err := preparePostgres(a.DBConnStr)
	if err != nil {
		return err
	}

	// initializing repositories
	ur, err := user.NewPostgres(pgpool, pgconf)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize postgres user repository at %s", a.DBConnStr)
	}

	ds := &service.DataStore{UserRepository: ur, BCryptCost: a.BCryptCost}

	var p []store.Privilege

	if err = json.Unmarshal([]byte(a.User.Privileges), &p); err != nil {
		return errors.Wrapf(err, "failed to unmarshal list of privileges %s", a.User.Privileges)
	}

	id, err := ds.AddUser(store.User{
		ID:         a.User.ID,
		Email:      a.User.Email,
		Privileges: p,
	}, a.User.Password)
	log.Printf("[INFO] registered user with id %s", id)
	return errors.Wrap(err, "failed to add user")
}

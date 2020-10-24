package main

// golangci-lint warns on the use of go-flags without alias
//noinspection GoRedundantImportAlias
import (
	"fmt"
	"os"

	log "github.com/go-pkgz/lgr"
	flags "github.com/jessevdk/go-flags"
	"github.com/yaattc/automatic-time-table-creation/backend/app/cmd"
)

// Opts describes cli arguments and flags to execute a command
type Opts struct {
	ServerCmd cmd.Server `command:"server"`

	Dbg bool `long:"dbg" env:"DEBUG" description:"turn on debug mode"`
}

var version = "unknown"

func main() {
	fmt.Printf("attc version: %s\n", version)
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)

	p.CommandHandler = func(command flags.Commander, args []string) error {
		setupLog(opts.Dbg)

		// commands implements CommonOptionsCommander to allow passing set of extra options defined for all commands
		c := command.(cmd.CommonOptionsCommander)
		c.SetCommon(cmd.CommonOpts{
			Version: version,
		})

		if err := command.Execute(args); err != nil {
			log.Printf("[ERROR] failed to execute command %+v", err)
		}
		return nil
	}

	// after failure command does not return non-zero code
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}

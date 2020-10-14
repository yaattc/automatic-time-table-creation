package main

// golangci-lint warns on the use of go-flags without alias
//noinspection GoRedundantImportAlias
import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	flags "github.com/jessevdk/go-flags"
	"github.com/yaattc/automatic-time-table-creation/backend/app/cmd"
)

// Opts describes cli arguments and flags to execute a command
type Opts struct {
	ServeCmd cmd.ServeCmd `command:"serve"`

	AttcURL string `long:"url" env:"ATTC_URL" required:"true" description:"url to attc"`
	Dbg     bool   `long:"dbg" env:"DEBUG" description:"turn on debug mode"`
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
			AttcURL: opts.AttcURL,
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
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: "INFO",
		Writer:   os.Stdout,
	}

	logFlags := log.Ldate | log.Ltime

	if dbg {
		logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile
		filter.MinLevel = "DEBUG"
	}

	log.SetFlags(logFlags)
	log.SetOutput(filter)
}

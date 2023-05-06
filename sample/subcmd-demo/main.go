package main

import (
	"fmt"

	"github.com/maargenton/go-cli"
	"gopkg.in/yaml.v3"
)

func main() {
	cli.Run(&cli.Command{
		// Handler:     &ServerCmd{},
		Description: "API server for demo service",

		SubCommands: []cli.Command{
			{
				Name:        "migrate",
				Description: "run DB migration",
				Handler:     &MigrateCmd{},
			},
		},
	})
}

type BaseOpts struct {
	Debug   bool `opts:"-d,--debug"   desc:"generating more diagnostic upon error"`
	Verbose bool `opts:"-v,--verbose" desc:"display additional information on startup"`

	DB       string `opts:"--db"        desc:"primary DB connection string"`
	DBPasswd string `opts:"--db-passwd" desc:"primary DB password"`
}

func (options *BaseOpts) Version() string {
	return "subcmd-demo v0.1.2"
}

// ---------------------------------------------------------------------------

type ServerCmd struct {
	BaseOpts
}

func (options *ServerCmd) Run() error {
	d, err := yaml.Marshal(options)
	if err != nil {
		return err
	}
	fmt.Printf("Running server with options:\n%v\n", string(d))
	return nil
}

// ---------------------------------------------------------------------------

type MigrateCmd struct {
	BaseOpts

	SkipDups bool `opts:"--skip-duplicates" desc:"skip duplicates during migration"`
}

func (options *MigrateCmd) Run() error {
	d, err := yaml.Marshal(options)
	if err != nil {
		return err
	}
	fmt.Printf("Running migration with options:\n%v\n", string(d))
	return nil
}

package main

import (
	"fmt"

	"github.com/maargenton/go-cli"
	"github.com/maargenton/go-cli/pkg/option"
	"gopkg.in/yaml.v3"
)

func main() {
	cli.Run(&cli.Command{
		Handler:     &sercatCmd{},
		Description: "Open a serial port and print all traffic to standard output",
	})
}

type sercatCmd struct {
	Port     *string `opts:"arg:1, name:port" desc:"name of the port to open"`
	Baudrate *uint32 `opts:"-b, --baudrate"   desc:"baudrate to use for communication"`
	Format   *string `opts:"-f, --format"     desc:"communcation format, e.g. 8N1"`

	Timestamp bool `opts:"-t,--timestamp" desc:"prefix every line with elapsed time"`
	Verbose   bool `opts:"-v,--verbose"   desc:"display additional information on startup"`
}

func (options *sercatCmd) Version() string {
	return "sercat-demo v0.1.2"
}

var baudrateCompletion = []string{
	"1200", "1800", "2400", "4800", "7200", "9600",
	"14400", "19200", "28800", "38400", "57600", "76800",
	"115200", "230400",
}

var formatCompletion = []string{
	"5N1", "6N1", "7N1", "8N1", "5N2", "6N2", "7N2", "8N2",
	"5O1", "6O1", "7O1", "8O1", "5O2", "6O2", "7O2", "8O2",
	"5E1", "6E1", "7E1", "8E1", "5E2", "6E2", "7E2", "8E2",
}

func (options *sercatCmd) Complete(opt *option.T, partial string) []string {
	if opt.Long == "baudrate" {
		return baudrateCompletion
	}
	if opt.Long == "format" {
		return formatCompletion
	}
	if opt.Position == 1 {
		return cli.MatchingFilenameCompletion(opt, "/dev/tty.*", partial)
	}
	return cli.DefaultCompletion(opt, partial)
}

func (options *sercatCmd) Run() error {

	// The following options use pointer type to allow getting fallback values
	// from a configuration file if not set. Manually set defaults here.
	if options.Baudrate == nil {
		var v uint32 = 115200
		options.Baudrate = &v
	}
	if options.Format == nil {
		var v string = "8N1"
		options.Format = &v
	}

	d, err := yaml.Marshal(options)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(d))
	return nil
}

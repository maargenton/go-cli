package main

import (
	"encoding/json"
	"fmt"

	"github.com/maargenton/go-cli"
	"github.com/maargenton/go-cli/pkg/option"
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
	return "sercat v1.0.0"
}

func (options *sercatCmd) Run() error {
	if options.Baudrate == nil {
		var v uint32 = 115200
		options.Baudrate = &v
	}
	if options.Format == nil {
		var v string = "8N1"
		options.Format = &v
	}

	d, err := json.Marshal(options)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(d))
	return nil
}

func (options *sercatCmd) Complete(opt *option.T, partial string) []option.Description {
	if opt.Long == "baudrate" {
		return baudrateCompletion
	}
	if opt.Long == "format" {
		return formatCompletion
	}
	if opt.Position == 1 {
		return cli.FilepathCompletion("/dev/tty.*", partial)
	}
	return cli.DefaultCompletion(partial)
}

var baudrateCompletion = []cli.Description{
	{Option: "1200"},
	{Option: "1800"},
	{Option: "2400"},
	{Option: "4800"},
	{Option: "7200"},
	{Option: "9600"},
	{Option: "14400"},
	{Option: "19200"},
	{Option: "28800"},
	{Option: "38400"},
	{Option: "57600"},
	{Option: "76800"},
	{Option: "115200"},
	{Option: "230400"},
}

var formatCompletion = []cli.Description{
	{Option: "5N1", Description: "5 bit data, no parity, 1 bit stop             "},
	{Option: "6N1", Description: "6 bit data, no parity, 1 bit stop             "},
	{Option: "7N1", Description: "7 bit data, no parity, 1 bit stop             "},
	{Option: "8N1", Description: "8 bit data, no parity, 1 bit stop             "},
	{Option: "5N2", Description: "5 bit data, no parity, 2 bit stop             "},
	{Option: "6N2", Description: "6 bit data, no parity, 2 bit stop             "},
	{Option: "7N2", Description: "7 bit data, no parity, 2 bit stop             "},
	{Option: "8N2", Description: "8 bit data, no parity, 2 bit stop             "},
	{Option: "5O1", Description: "5 bit data, odd parity, 1 bit stop            "},
	{Option: "6O1", Description: "6 bit data, odd parity, 1 bit stop            "},
	{Option: "7O1", Description: "7 bit data, odd parity, 1 bit stop            "},
	{Option: "8O1", Description: "8 bit data, odd parity, 1 bit stop            "},
	{Option: "5O2", Description: "5 bit data, odd parity, 2 bit stop            "},
	{Option: "6O2", Description: "6 bit data, odd parity, 2 bit stop            "},
	{Option: "7O2", Description: "7 bit data, odd parity, 2 bit stop            "},
	{Option: "8O2", Description: "8 bit data, odd parity, 2 bit stop            "},
	{Option: "5E1", Description: "5 bit data, even parity, 1 bit stop           "},
	{Option: "6E1", Description: "6 bit data, even parity, 1 bit stop           "},
	{Option: "7E1", Description: "7 bit data, even parity, 1 bit stop           "},
	{Option: "8E1", Description: "8 bit data, even parity, 1 bit stop           "},
	{Option: "5E2", Description: "5 bit data, even parity, 2 bit stop           "},
	{Option: "6E2", Description: "6 bit data, even parity, 2 bit stop           "},
	{Option: "7E2", Description: "7 bit data, even parity, 2 bit stop           "},
	{Option: "8E2", Description: "8 bit data, even parity, 2 bit stop           "},
}

package main

import (
	"github.com/maargenton/go-cli"
	"github.com/maargenton/go-cli/pkg/enumer"
)

func main() {
	cli.Run(&cli.Command{
		Handler:     &enumer.Cmd{},
		Description: "Generate flag.Value interface for enumerated types",
	})
}

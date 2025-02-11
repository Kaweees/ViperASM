package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type argsParsed struct {
	args
}

type args struct {
	// File config
	FileName string `arg:"required" help:"File to assemble"`
	// Logging config
	Logging bool `arg:"-l,--logging" help:"Enable logging"`
}

// Returns a human-readable version string
func (args) Version() string {
	return fmt.Sprintf("Version: %v, commit: %v, built at: %v", version, commit, date)
}

// Returns a description of the program
func (args) Description() string {
	return "A simple assembler for the MIPS architecture"
}

// Returns the parsed CLI arguments
func GetCliArgs() (cli argsParsed, err error) {
	rawCli := args{}
	rawCli.Logging = false

	arg.MustParse(&rawCli)
	cli.args = rawCli

	return cli, nil
}

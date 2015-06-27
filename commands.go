package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/spiritloose/consulkv/command"
)

// Commands is the mapping of all the available consulkv commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	Commands = map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &command.ListCommand{UI: ui}, nil
		},
		"cat": func() (cli.Command, error) {
			return &command.CatCommand{UI: ui, Output: os.Stdout}, nil
		},
		"put": func() (cli.Command, error) {
			return &command.PutCommand{UI: ui}, nil
		},
		"edit": func() (cli.Command, error) {
			return &command.EditCommand{UI: ui}, nil
		},
		"delete": func() (cli.Command, error) {
			return &command.DeleteCommand{UI: ui}, nil
		},
		"flags": func() (cli.Command, error) {
			return &command.FlagsCommand{UI: ui}, nil
		},
		"dump": func() (cli.Command, error) {
			return &command.DumpCommand{UI: ui, Output: os.Stdout}, nil
		},
		"load": func() (cli.Command, error) {
			return &command.LoadCommand{UI: ui, Input: os.Stdin}, nil
		},
	}
}

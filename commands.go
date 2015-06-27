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
			return &command.List{UI: ui}, nil
		},
		"cat": func() (cli.Command, error) {
			return &command.Cat{UI: ui, Output: os.Stdout}, nil
		},
		"put": func() (cli.Command, error) {
			return &command.Put{UI: ui}, nil
		},
		"edit": func() (cli.Command, error) {
			return &command.Edit{UI: ui}, nil
		},
		"delete": func() (cli.Command, error) {
			return &command.Delete{UI: ui}, nil
		},
		"flags": func() (cli.Command, error) {
			return &command.Flags{UI: ui}, nil
		},
		"dump": func() (cli.Command, error) {
			return &command.Dump{UI: ui, Output: os.Stdout}, nil
		},
		"load": func() (cli.Command, error) {
			return &command.Load{UI: ui, Input: os.Stdin}, nil
		},
	}
}

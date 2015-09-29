package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

// Version string of consulkv
var VERSION = "0.2.0-dev"

func main() {
	c := cli.NewCLI("consulkv", VERSION)
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "-v" || arg == "--version" {
			fmt.Println(VERSION)
			os.Exit(0)
		}
	}

	c.Args = args
	c.Commands = Commands

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitStatus)
}

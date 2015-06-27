package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

// Version string of consulkv
var VERSION = "HEAD"

func main() {
	c := cli.NewCLI("consulkv", VERSION)
	c.Args = os.Args[1:]
	c.Commands = Commands

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitStatus)
}

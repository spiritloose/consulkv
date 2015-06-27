package command

import (
	"flag"
	"fmt"
	"io"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

// CatCommand is a Command implementation that concatenate and print values.
type CatCommand struct {
	UI     cli.Ui
	Output io.Writer
}

// Help prints the Help text for the cat command.
func (c *CatCommand) Help() string {
	return "Usage: consulkv cat [-datacenter=] KEY..."
}

// Synopsis provides a precis of the cat command.
func (c *CatCommand) Synopsis() string {
	return "Concatenate and print values"
}

// Run runs the cat command.
func (c *CatCommand) Run(args []string) int {
	var datacenter string
	cmdFlags := flag.NewFlagSet("cat", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&datacenter, "datacenter", "", "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	args = cmdFlags.Args()
	if len(args) == 0 {
		c.UI.Error("Key must be specified")
		return 1
	}
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error connecting to Consul agent: %s", err))
		return 1
	}
	kv := client.KV()

	failed := false
	for _, key := range args {
		pair, _, err := kv.Get(key, &api.QueryOptions{Datacenter: datacenter})
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error getting key: %s", err))
			return 1
		}
		if pair == nil {
			c.UI.Error(fmt.Sprintf("cat: %s: No such key", key))
			failed = true
		} else {
			c.Output.Write(pair.Value)
		}
	}

	if failed {
		return 1
	}
	return 0
}

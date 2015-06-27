package command

import (
	"flag"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

// DeleteCommand is a Command implementation that removes key/value entries.
type DeleteCommand struct {
	UI cli.Ui
}

// Help prints the Help text for the delete command.
func (c *DeleteCommand) Help() string {
	return "Usage: consulkv delete [-datacenter=] [-recursive] KEY..."
}

// Synopsis provides a precis of the delete command.
func (c *DeleteCommand) Synopsis() string {
	return "Remove key/value entries"
}

// Run runs the delete command.
func (c *DeleteCommand) Run(args []string) int {
	var datacenter string
	var recursive bool
	cmdFlags := flag.NewFlagSet("delete", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&datacenter, "datacenter", "", "")
	cmdFlags.BoolVar(&recursive, "recursive", false, "")
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
	}

	kv := client.KV()
	options := api.WriteOptions{Datacenter: datacenter}
	for _, key := range args {
		if recursive {
			_, err = kv.DeleteTree(key, &options)
		} else {
			_, err = kv.Delete(key, &options)
		}
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error deleting key: %s", err))
			return 1
		}
	}

	return 0
}

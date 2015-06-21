package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type Cat struct {
	Ui cli.Ui
}

func (c *Cat) Help() string {
	return "Usage: consulkv cat [-datacenter=] KEY..."
}

func (c *Cat) Synopsis() string {
	return "Concatenate and print values"
}

func (c *Cat) Run(args []string) int {
	var datacenter string
	cmdFlags := flag.NewFlagSet("cat", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&datacenter, "datacenter", "", "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	args = cmdFlags.Args()
	if len(args) == 0 {
		c.Ui.Error("Key must be specified")
		return 1
	}
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error connecting to Consul agent: %s", err))
		return 1
	}
	kv := client.KV()

	failed := false
	for _, key := range args {
		pair, _, err := kv.Get(key, &api.QueryOptions{Datacenter: datacenter})
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error getting key: %s", err))
			return 1
		}
		if pair == nil {
			c.Ui.Error(fmt.Sprintf("cat: %s: No such key", key))
			failed = true
		} else {
			os.Stdout.Write(pair.Value)
		}
	}

	if failed {
		return 1
	}
	return 0
}

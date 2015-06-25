package command

import (
	"flag"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type List struct {
	UI cli.Ui
}

func (c *List) Help() string {
	return "Usage: consulkv list [-datacenter=] [-separator=/] [PREFIX...]"
}

func (c *List) Synopsis() string {
	return "List keys"
}

func (c *List) Run(args []string) int {
	var datacenter string
	var separator string
	cmdFlags := flag.NewFlagSet("list", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&datacenter, "datacenter", "", "")
	cmdFlags.StringVar(&separator, "separator", "/", "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	args = cmdFlags.Args()
	if len(args) == 0 {
		args = append(args, "")
	}
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error connecting to Consul agent: %s", err))
		return 1
	}
	kv := client.KV()

	options := api.QueryOptions{Datacenter: datacenter}
	for _, prefix := range args {
		keys, _, err := kv.Keys(prefix, separator, &options)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error getting keys: %s", err))
			return 1
		}
		for _, key := range keys {
			c.UI.Output(key)
		}
	}

	return 0
}

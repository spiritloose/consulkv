package command

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type Flags struct {
	UI cli.Ui
}

func (c *Flags) Help() string {
	return "Usage: consulkv flags [-datacenter=] KEY [FLAGS]"
}

func (c *Flags) Synopsis() string {
	return "Get/Set flags"
}

func (c *Flags) Run(args []string) int {
	var datacenter string
	cmdFlags := flag.NewFlagSet("flags", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&datacenter, "datacenter", "", "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	args = cmdFlags.Args()

	argsLen := len(args)
	var key string
	var flags uint64
	var err error
	switch argsLen {
	case 0:
		c.UI.Error("Key must be specified")
		return 1
	case 1:
		key = args[0]
	case 2:
		key = args[0]
		flags, err = strconv.ParseUint(args[1], 10, 0)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Invalid flags format: %s", err))
			return 1
		}
	default:
		c.UI.Error("Too many arguments")
		return 1
	}

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error connecting to Consul agent: %s", err))
		return 1
	}
	kv := client.KV()
	pair, _, err := kv.Get(key, &api.QueryOptions{Datacenter: datacenter})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting key: %s", err))
		return 1
	}
	if pair == nil {
		c.UI.Error(fmt.Sprintf("flags: %s: No such key", key))
		return 1
	}

	if argsLen == 1 {
		c.UI.Output(strconv.FormatUint(pair.Flags, 10))
		return 0
	}

	pair.Flags = flags
	success, _, err := kv.CAS(pair, &api.WriteOptions{Datacenter: datacenter})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting key/value: %s", err))
		return 1
	}
	if !success {
		c.UI.Error("The key has been changed since reading it")
		return 1
	}

	return 0
}

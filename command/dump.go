package command

import (
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type Dump struct {
	UI     cli.Ui
	Output io.Writer
}

func (c *Dump) Help() string {
	return "Usage: consulkv dump [-datacenter=] [PREFIX...]"
}

func (c *Dump) Synopsis() string {
	return "Dump key/values"
}

func (c *Dump) Run(args []string) int {
	var datacenter string
	cmdFlags := flag.NewFlagSet("list", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&datacenter, "datacenter", "", "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	args = cmdFlags.Args()
	if len(args) == 0 {
		args = []string{"/"}
	}
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error connecting to Consul agent: %s", err))
		return 1
	}
	kv := client.KV()

	options := api.QueryOptions{Datacenter: datacenter}
	writer := csv.NewWriter(c.Output)
	writer.Comma = '\t'
	for _, prefix := range args {
		pairs, _, err := kv.List(prefix, &options)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error listing key/value: %s", err))
			return 1
		}
		for _, pair := range pairs {
			value := base64.StdEncoding.EncodeToString(pair.Value)
			writer.Write([]string{pair.Key, value, strconv.FormatUint(pair.Flags, 10)})
		}
	}
	writer.Flush()

	return 0
}

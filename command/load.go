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

type LoadCommand struct {
	UI    cli.Ui
	Input io.Reader
}

func (c *LoadCommand) Help() string {
	return "Usage: consulkv load [-datacenter=]"
}

func (c *LoadCommand) Synopsis() string {
	return "Load key/values"
}

func (c *LoadCommand) Run(args []string) int {
	var datacenter string
	cmdFlags := flag.NewFlagSet("load", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&datacenter, "datacenter", "", "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error connecting to Consul agent: %s", err))
		return 1
	}
	kv := client.KV()

	options := api.WriteOptions{Datacenter: datacenter}
	reader := csv.NewReader(c.Input)
	reader.Comma = '\t'
	reader.FieldsPerRecord = 3
	reader.LazyQuotes = true
	lineNum := 0
	failed := false
	for {
		lineNum++
		columns, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			c.UI.Error(err.Error())
			failed = true
			continue
		}
		key := columns[0]
		if len(key) == 0 {
			c.UI.Error(fmt.Sprintf("Line %d: Key must be specified", lineNum))
			failed = true
			continue
		}

		valueStr := columns[1]
		var value []byte
		if len(valueStr) > 0 {
			value, err = base64.StdEncoding.DecodeString(valueStr)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Line %d: Error decoding value err: %s", lineNum, err))
				failed = true
				continue
			}
		}

		flagsStr := columns[2]
		if len(flagsStr) == 0 {
			c.UI.Error(fmt.Sprintf("Line %d: Flags must be specified err: %s", lineNum, err))
			failed = true
			continue
		}
		var flags uint64
		flags, err = strconv.ParseUint(flagsStr, 10, 0)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Line %d: Error parsing flags err: %s", lineNum, err))
			failed = true
			continue
		}

		pair := api.KVPair{Key: key, Value: value, Flags: flags}
		_, err = kv.Put(&pair, &options)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Line %d: Error putting Key/value err: %s", lineNum, err))
			failed = true
			continue
		}
	}

	if failed {
		return 1
	}
	return 0
}

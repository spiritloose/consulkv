package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/cli"
)

type EditCommand struct {
	UI cli.Ui
}

func (c *EditCommand) Help() string {
	return "Usage: consulkv edit [-datacenter=] [-flags=0] [-chomp] KEY"
}

func (c *EditCommand) Synopsis() string {
	return "Edit value using a editor"
}

func (c *EditCommand) Run(args []string) int {
	var datacenter string
	var flags uint64
	var chomp bool
	cmdFlags := flag.NewFlagSet("edit", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&datacenter, "datacenter", "", "")
	cmdFlags.Uint64Var(&flags, "flags", 0, "")
	cmdFlags.BoolVar(&chomp, "chomp", false, "")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	args = cmdFlags.Args()
	if len(args) != 1 {
		c.UI.Error("Key must be specified")
		return 1
	}
	key := args[0]
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

	file, err := ioutil.TempFile(os.TempDir(), "consulkv-")
	defer os.Remove(file.Name())
	if pair != nil {
		file.Write(pair.Value)
	}

	fi, err := file.Stat()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting stats of tempfile: %s", err))
		return 1
	}
	beforeTime := fi.ModTime()

	err = file.Close()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error closing file: %s", err))
		return 1
	}

	err = c.execEditor(file.Name())
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error launching editor: %s", err))
		return 1
	}

	file, err = os.Open(file.Name())
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error opening file: %s", err))
		return 1
	}
	fi, err = file.Stat()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting stats of tempfile: %s", err))
		return 1
	}

	if beforeTime == fi.ModTime() {
		c.UI.Warn("Not modified. Aborted")
		return 1
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading file: %s", err))
		return 1
	}
	var value []byte
	if chomp {
		value = []byte(strings.TrimRight(string(content), "\n"))
	} else {
		value = content
	}

	writeOptions := api.WriteOptions{Datacenter: datacenter}
	newPair := api.KVPair{Key: key, Value: value, Flags: flags}
	if pair != nil {
		if newPair.Flags == 0 {
			newPair.Flags = pair.Flags
		}
		newPair.ModifyIndex = pair.ModifyIndex
		success, _, err := kv.CAS(&newPair, &writeOptions)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error setting key/value: %s", err))
			return 1
		}
		if success {
			return 0
		}
		success, err = c.askOverwrite()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error reading the answer from console: %s", err))
			return 1
		}
		if !success {
			c.UI.Error("Aborted")
			return 1
		}
	}

	_, err = kv.Put(&newPair, &writeOptions)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting key/value: %s", err))
		return 1
	}

	return 0
}

func (c *EditCommand) execEditor(filename string) error {
	editor := c.getEditor()
	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (c *EditCommand) getEditor() string {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vim"
	}
	return editor
}

func (c *EditCommand) askOverwrite() (bool, error) {
	result, err := c.UI.Ask(`WARNING: The key has been changed since reading it!!!
Do you really want to write to it (y/n)?`)
	if err != nil {
		return false, err
	}
	return result == "y", nil
}

package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestEditCommand(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	srv.SetKV("foo", []byte("bar"))

	ui := new(cli.MockUi)
	c := &EditCommand{UI: ui}

	file, _ := ioutil.TempFile(os.TempDir(), "consulkv-test")
	defer os.Remove(file.Name())
	fmt.Fprint(file, strings.TrimSpace(`
#!/bin/sh
echo barbar > $1
	`))

	os.Setenv("CONSUL_HTTP_ADDR", srv.HTTPAddr)
	os.Setenv("EDITOR", fmt.Sprintf("sh %s", file.Name()))
	args := []string{"foo"}
	code := c.Run(args)
	if code != 0 {
		t.Fatalf("Unexpected code: %d err: %s", code, ui.ErrorWriter.String())
	}
	val := srv.GetKV("foo")
	if string(val) != "barbar\n" {
		t.Fatalf("Invalid value %s", val)
	}
}

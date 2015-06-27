package command

import (
	"bytes"
	"os"
	"testing"

	"github.com/hashicorp/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestDumpCommand(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	srv.SetKV("foo", []byte("bar"))
	srv.SetKV("bar", []byte("baz"))

	ui := new(cli.MockUi)
	var output bytes.Buffer
	c := &DumpCommand{UI: ui, Output: &output}

	os.Setenv("CONSUL_HTTP_ADDR", srv.HTTPAddr)
	args := []string{}
	code := c.Run(args)
	if code != 0 {
		t.Fatalf("Unexpected code: %d err: %s", code, ui.ErrorWriter.String())
	}

	str := output.String()
	if str != "bar\tYmF6\t0\nfoo\tYmFy\t0\n" {
		t.Fatalf("Unexpected data: %s", str)
	}
}

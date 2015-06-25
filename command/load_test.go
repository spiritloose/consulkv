package command

import (
	"bytes"
	"os"
	"testing"

	"github.com/hashicorp/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestLoadCommand(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	ui := new(cli.MockUi)
	input := bytes.NewBufferString("bar\tYmF6\t0\nfoo\tYmFy\t0\n")

	c := &Load{UI: ui, Input: input}

	os.Setenv("CONSUL_HTTP_ADDR", srv.HTTPAddr)
	args := []string{}
	code := c.Run(args)
	if code != 0 {
		t.Fatalf("Unexpected code: %d err: %s", code, ui.ErrorWriter.String())
	}

	res := string(srv.GetKV("foo"))
	if res != "bar" {
		t.Fatalf("Unexpected string: %s err: %s", res, ui.ErrorWriter.String())
	}
	res = string(srv.GetKV("bar"))
	if res != "baz" {
		t.Fatalf("Unexpected string: %s err: %s", res, ui.ErrorWriter.String())
	}
}

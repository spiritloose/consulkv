package command

import (
	"bytes"
	"os"
	"testing"

	"github.com/hashicorp/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestPutCommand(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	ui := new(cli.MockUi)
	c := &PutCommand{UI: ui}

	os.Setenv("CONSUL_HTTP_ADDR", srv.HTTPAddr)
	args := []string{"foo", "bar"}
	code := c.Run(args)
	if code != 0 {
		t.Fatalf("Unexpected code: %d err: %s", code, ui.ErrorWriter.String())
	}
	val := srv.GetKV("foo")
	if string(val) != "bar" {
		t.Fatalf("Invalid value %s", val)
	}
}

func TestPutCommandStdin(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	ui := new(cli.MockUi)
	input := bytes.NewBufferString("bar")
	c := &PutCommand{UI: ui, Input: input}

	os.Setenv("CONSUL_HTTP_ADDR", srv.HTTPAddr)
	args := []string{"foo"}
	code := c.Run(args)
	if code != 0 {
		t.Fatalf("Unexpected code: %d err: %s", code, ui.ErrorWriter.String())
	}
	val := srv.GetKV("foo")
	if string(val) != "bar" {
		t.Fatalf("Invalid value %s", val)
	}
}

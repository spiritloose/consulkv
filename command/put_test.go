package command

import (
	"os"
	"testing"

	"github.com/hashicorp/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestPutCommand(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	ui := new(cli.MockUi)
	c := &Put{UI: ui}

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

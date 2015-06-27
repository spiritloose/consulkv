package command

import (
	"os"
	"testing"

	"github.com/hashicorp/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestFlagsCommand(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	srv.SetKV("foo", []byte("bar"))

	ui := new(cli.MockUi)
	c := &FlagsCommand{UI: ui}

	os.Setenv("CONSUL_HTTP_ADDR", srv.HTTPAddr)
	args := []string{"foo"}
	code := c.Run(args)
	if code != 0 {
		t.Fatalf("Unexpected code: %d err: %s", code, ui.ErrorWriter.String())
	}

	str := ui.OutputWriter.String()
	if str != "0\n" {
		t.Fatalf("Unexpected data: %s", str)
	}
}

package command

import (
	"os"
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

	os.Setenv("CONSUL_HTTP_ADDR", srv.HTTPAddr)
	os.Setenv("EDITOR", "sed -i '' -e 's/bar/barbar/'")
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

package command

import (
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestListCommand(t *testing.T) {
	srv := testutil.NewTestServer(t)
	defer srv.Stop()

	srv.SetKV("foo", []byte("bar"))
	srv.SetKV("bar", []byte("baz"))

	ui := new(cli.MockUi)
	c := &List{UI: ui}

	os.Setenv("CONSUL_HTTP_ADDR", srv.HTTPAddr)
	args := []string{}
	code := c.Run(args)
	if code != 0 {
		t.Fatalf("Unexpected code: %d err: %s", code, ui.ErrorWriter.String())
	}

	keys := strings.Split(strings.TrimRight(ui.OutputWriter.String(), "\n"), "\n")
	if len(keys) != 2 {
		t.Fatalf("Unexpected key size: %d datq: %v", len(keys), keys)
	}
	sort.Strings(keys)
	if !reflect.DeepEqual(keys, []string{"bar", "foo"}) {
		t.Fatalf("Unexpected data: %v", keys)
	}
}

package main

import (
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/gliderlabs/hostctl/providers"
)

func TestIpCmd(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test1",
		IP:   "127.0.0.1",
	})

	stdout, stderr := testRunCmd(t, "hostctl ip test1", 0, provider, nil)
	ensure.DeepEqual(t, stdout.String(), "127.0.0.1\n")
	ensure.DeepEqual(t, stderr.String(), "")
}

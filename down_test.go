package main

import (
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/gliderlabs/hostctl/providers"
)

func TestDownCmd(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test1",
	})
	provider.Create(providers.Host{
		Name: "test2",
	})

	stdout, stderr := testRunCmd(t, "hostctl down test1", 0, provider, nil)
	ensure.DeepEqual(t, stdout.String(), "")
	ensure.DeepEqual(t, stderr.String(), "\n")

	ensure.DeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
	ensure.NotDeepEqual(t, provider.Get("test2"), (*providers.Host)(nil))
}

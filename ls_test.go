package main

import (
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/gliderlabs/hostctl/providers"
)

func TestListBasic(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test1",
	})
	provider.Create(providers.Host{
		Name: "test2",
	})

	stdout, stderr := testRunCmd(t, "hostctl ls", 0, provider, nil)
	ensure.DeepEqual(t, stdout.String(), "test1\ntest2\n")
	ensure.DeepEqual(t, stderr.String(), "")
}

func TestListPatternSingle(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test1",
	})
	provider.Create(providers.Host{
		Name: "test2",
	})

	stdout, stderr := testRunCmd(t, "hostctl ls test1", 0, provider, nil)
	ensure.DeepEqual(t, stdout.String(), "test1\n")
	ensure.DeepEqual(t, stderr.String(), "")
}

func TestListPatternPartial(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test1",
	})
	provider.Create(providers.Host{
		Name: "test2",
	})
	provider.Create(providers.Host{
		Name: "foobar",
	})

	stdout, stderr := testRunCmd(t, "hostctl ls te*", 0, provider, nil)
	ensure.DeepEqual(t, stdout.String(), "test1\ntest2\n")
	ensure.DeepEqual(t, stderr.String(), "")
}

package main

import (
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/gliderlabs/hostctl/providers"
)

func TestUpSimple(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)

	stdout, stderr := testRunCmd(t, "hostctl up test1", 0, provider, nil)
	ensure.DeepEqual(t, stderr.String(), "\n")
	ensure.DeepEqual(t, stdout.String(), "")

	ensure.NotDeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
}

func TestUpExists(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test1",
	})

	stdout, stderr := testRunCmd(t, "hostctl up test1", 0, provider, nil)
	ensure.DeepEqual(t, stderr.String(), "\n")
	ensure.DeepEqual(t, stdout.String(), "")

	ensure.NotDeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
}

func TestUpMultiple(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)

	stdout, stderr := testRunCmd(t, "hostctl up test1 test2 test3", 0, provider, nil)
	ensure.DeepEqual(t, stderr.String(), "\n")
	ensure.DeepEqual(t, stdout.String(), "")

	ensure.NotDeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
	ensure.NotDeepEqual(t, provider.Get("test2"), (*providers.Host)(nil))
	ensure.NotDeepEqual(t, provider.Get("test3"), (*providers.Host)(nil))
}

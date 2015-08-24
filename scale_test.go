package main

import (
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/gliderlabs/hostctl/providers"
)

func TestScaleUpFromZero(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)

	stdout, stderr := testRunCmd(t, "hostctl scale test 3", 0, provider, nil)
	ensure.DeepEqual(t, stderr.String(), "\n")
	ensure.DeepEqual(t, stdout.String(), "")

	ensure.NotDeepEqual(t, provider.Get("test.0"), (*providers.Host)(nil))
	ensure.NotDeepEqual(t, provider.Get("test.1"), (*providers.Host)(nil))
	ensure.NotDeepEqual(t, provider.Get("test.2"), (*providers.Host)(nil))
}

func TestScaleDownToZero(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test.0",
	})
	provider.Create(providers.Host{
		Name: "test.1",
	})
	provider.Create(providers.Host{
		Name: "test.2",
	})

	stdout, stderr := testRunCmd(t, "hostctl scale test 0", 0, provider, nil)
	ensure.DeepEqual(t, stderr.String(), "\n")
	ensure.DeepEqual(t, stdout.String(), "")

	ensure.DeepEqual(t, provider.Get("test.2"), (*providers.Host)(nil))
	ensure.DeepEqual(t, provider.Get("test.1"), (*providers.Host)(nil))
	ensure.DeepEqual(t, provider.Get("test.0"), (*providers.Host)(nil))
}

func TestScaleDownToOne(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test.0",
	})
	provider.Create(providers.Host{
		Name: "test.1",
	})
	provider.Create(providers.Host{
		Name: "test.2",
	})

	stdout, stderr := testRunCmd(t, "hostctl scale test 1", 0, provider, nil)
	ensure.DeepEqual(t, stderr.String(), "\n")
	ensure.DeepEqual(t, stdout.String(), "")

	ensure.NotDeepEqual(t, provider.Get("test.0"), (*providers.Host)(nil))
	ensure.DeepEqual(t, provider.Get("test.1"), (*providers.Host)(nil))
	ensure.DeepEqual(t, provider.Get("test.2"), (*providers.Host)(nil))
}

func TestScaleUpFromOne(t *testing.T) {
	t.Parallel()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test.0",
	})

	stdout, stderr := testRunCmd(t, "hostctl scale test 3", 0, provider, nil)
	ensure.DeepEqual(t, stderr.String(), "\n")
	ensure.DeepEqual(t, stdout.String(), "")

	ensure.NotDeepEqual(t, provider.Get("test.0"), (*providers.Host)(nil))
	ensure.NotDeepEqual(t, provider.Get("test.1"), (*providers.Host)(nil))
	ensure.NotDeepEqual(t, provider.Get("test.2"), (*providers.Host)(nil))
}

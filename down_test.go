package main

import (
	"bytes"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/gliderlabs/hostctl/providers"
)

func TestDownCmd(t *testing.T) {
	resetEnv()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test1",
	})
	provider.Create(providers.Host{
		Name: "test2",
	})
	providerName = "test"
	providers.Register(provider, providerName)
	var out, err bytes.Buffer
	status, exit := captureExit()
	runCmd(downCmd, &Context{
		Out:  &out,
		Err:  &err,
		Args: []string{"test1"},
		Exit: exit,
	})
	ensure.DeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
	ensure.NotDeepEqual(t, provider.Get("test2"), (*providers.Host)(nil))
	ensure.DeepEqual(t, err.String(), "\n")
	ensure.DeepEqual(t, status(), 0)
}

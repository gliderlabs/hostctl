package main

import (
	"bytes"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/gliderlabs/hostctl/providers"
)

func TestUpSimple(t *testing.T) {
	resetEnv()
	provider := new(providers.TestProvider)
	providerName = "test"
	providers.Register(provider, providerName)
	ensure.DeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
	var out, err bytes.Buffer
	status, exit := captureExit()
	runCmd(upCmd, &Context{
		Out:  &out,
		Err:  &err,
		Args: []string{"test1"},
		Exit: exit,
	})
	ensure.NotDeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
	ensure.DeepEqual(t, err.String(), "\n")
	ensure.DeepEqual(t, status(), 0)
}

func TestUpExists(t *testing.T) {
	resetEnv()
	provider := new(providers.TestProvider)
	provider.Create(providers.Host{
		Name: "test1",
	})
	providerName = "test"
	providers.Register(provider, providerName)
	ensure.NotDeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
	var out, err bytes.Buffer
	status, exit := captureExit()
	runCmd(upCmd, &Context{
		Out:  &out,
		Err:  &err,
		Args: []string{"test1"},
		Exit: exit,
	})
	ensure.NotDeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
	ensure.DeepEqual(t, err.String(), "\n")
	ensure.DeepEqual(t, status(), 0)
}

func TestUpMultiple(t *testing.T) {
	resetEnv()
	provider := new(providers.TestProvider)
	providerName = "test"
	providers.Register(provider, providerName)
	ensure.DeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
	ensure.DeepEqual(t, provider.Get("test2"), (*providers.Host)(nil))
	ensure.DeepEqual(t, provider.Get("test3"), (*providers.Host)(nil))
	var out, err bytes.Buffer
	status, exit := captureExit()
	runCmd(upCmd, &Context{
		Out:  &out,
		Err:  &err,
		Args: []string{"test1", "test2", "test3"},
		Exit: exit,
	})
	ensure.NotDeepEqual(t, provider.Get("test1"), (*providers.Host)(nil))
	ensure.NotDeepEqual(t, provider.Get("test2"), (*providers.Host)(nil))
	ensure.NotDeepEqual(t, provider.Get("test3"), (*providers.Host)(nil))
	ensure.DeepEqual(t, err.String(), "\n")
	ensure.DeepEqual(t, status(), 0)
}

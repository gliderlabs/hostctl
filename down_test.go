package main

import (
	"bytes"
	"fmt"
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
	runCmd(downCmd, &Context{
		Out:  &out,
		Err:  &err,
		Args: []string{"test1"},
	})
	fmt.Printf("%v", provider.Get("test1"))
	if provider.Get("test1") == nil {
		fmt.Println("its nil")
	}
	ensure.Nil(t, provider.Get("test1"))

	//ensure.NotNil(t, provider.Get("test2"))
	ensure.DeepEqual(t, err.String(), "")
}

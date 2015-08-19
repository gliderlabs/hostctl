package main

import (
	"bytes"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestEnvCmd(t *testing.T) {
	resetEnv()
	var out, err bytes.Buffer
	runCmd(envCmd, &Context{
		Out: &out,
		Err: &err,
	})
	ensure.StringContains(t, out.String(), "HOSTCTL_PROVIDER=")
	ensure.StringContains(t, out.String(), "HOSTCTL_NAMESPACE=")
	ensure.DeepEqual(t, err.String(), "")
}

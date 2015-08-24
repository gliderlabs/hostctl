package main

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestEnvCmd(t *testing.T) {
	t.Parallel()
	stdout, stderr := testRunCmd(t, "hostctl env", 0, nil, nil)
	ensure.StringContains(t, stdout.String(), "HOSTCTL_PROVIDER=")
	ensure.StringContains(t, stdout.String(), "HOSTCTL_NAMESPACE=")
	ensure.DeepEqual(t, stderr.String(), "")
}

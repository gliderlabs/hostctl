package main

import (
	"bytes"
	"testing"

	"github.com/facebookgo/ensure"
)

func runCmd(cmd *Command, context *Context) {
	context.Cmd = cmd.setup().Cmd
	context.Cmd.SetOutput(context)
	cmd.Run(context)
}

func resetEnv() {
	providerName = ""
	defaultName = ""
	namespace = ""
	hostImage = ""
	hostFlavor = ""
	hostRegion = ""
	hostKeyname = ""
	hostUserdata = ""
	user = ""
}

func captureExit() (func() int, func(int)) {
	ch := make(chan int, 1)
	return func() int {
			select {
			case status := <-ch:
				return status
			default:
				return 0
			}
		}, func(status int) {
			ch <- status
		}
}

func TestHostctlCmd(t *testing.T) {
	resetEnv()
	var out, err bytes.Buffer
	runCmd(Hostctl, &Context{
		Out: &out,
		Err: &err,
	})
	ensure.StringContains(t, out.String(), Hostctl.Short)
	ensure.StringContains(t, out.String(), "Usage:")
	ensure.StringContains(t, out.String(), "Available Commands:")
	ensure.StringContains(t, out.String(), "Flags:")
	ensure.DeepEqual(t, err.String(), "")
}

func TestVersionCmd(t *testing.T) {
	resetEnv()
	Version = "0"
	var out, err bytes.Buffer
	runCmd(versionCmd, &Context{
		Out: &out,
		Err: &err,
	})
	ensure.DeepEqual(t, out.String(), "0\n")
	ensure.DeepEqual(t, err.String(), "")
}

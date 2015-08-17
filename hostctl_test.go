package main

import (
	"bytes"
	"testing"

	"github.com/facebookgo/ensure"
)

func runCmd(cmd *Command, context *Context) {
	context.Cmd = cmd.setup().Cmd
	cmd.Run(context)
}

func TestVersionCmd(t *testing.T) {
	Version = "0"
	var out, err bytes.Buffer
	runCmd(cmdVersion, &Context{
		Out: &out,
		Err: &err,
	})
	ensure.DeepEqual(t, out.String(), "0\n")
	ensure.DeepEqual(t, err.String(), "")
}

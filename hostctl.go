package main

import (
	"fmt"
	"os"

	"github.com/MattAitchison/env"
)

var (
	Version string

	providerName = env.String("HOSTCTL_PROVIDER", "digitalocean", "cloud provider")
	defaultName  = env.String("HOSTCTL_NAME", "", "optional default name")
	namespace    = env.String("HOSTCTL_NAMESPACE", "", "optional namespace for names")

	hostImage    = env.String("HOSTCTL_IMAGE", "", "vm image")
	hostFlavor   = env.String("HOSTCTL_FLAVOR", "", "vm flavor")
	hostRegion   = env.String("HOSTCTL_REGION", "", "vm region")
	hostKeyname  = env.String("HOSTCTL_KEYNAME", "", "vm keyname")
	hostUserdata = env.String("HOSTCTL_USERDATA", "", "vm user data")

	user = env.String("HOSTCTL_USER", os.Getenv("USER"), "ssh user")
)

func main() {
	Hostctl.Context = Context{
		Out:  os.Stdout,
		Err:  os.Stderr,
		In:   os.Stdin,
		Exit: os.Exit,
	}
	Hostctl.AddCommand(cmdVersion)
	fatal(Hostctl.Execute())
}

var Hostctl = &Command{
	Use:   "hostctl",
	Short: "An opinionated tool for provisioning cloud VMs",
	Run: func(ctx *Context) {
		ctx.Cmd.Help()
	},
}

var cmdVersion = &Command{
	Use:   "version",
	Short: "Show version",
	Run: func(ctx *Context) {
		fmt.Fprintln(ctx, Version)
	},
}

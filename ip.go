package main

import (
	"fmt"
	"os"

	"github.com/gliderlabs/hostctl/providers"
)

func init() {
	Hostctl.AddCommand(ipCmd)
}

var ipCmd = &Command{
	Use:   "ip <name>",
	Short: "Show IP for host",
	Run: func(ctx *Context) {
		args := ctx.Args
		if len(args) < 1 && defaultName == "" {
			ctx.Cmd.Usage()
			os.Exit(1)
		}
		provider, err := providers.Get(providerName, true)
		fatal(err)
		host := provider.Get(namespace + optArg(args, 0, defaultName))
		if host == nil {
			os.Exit(1)
		}
		fmt.Println(host.IP)
	},
}

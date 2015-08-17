package main

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/hostctl/providers"
)

var fullNames bool

func init() {
	listCmd.Flags().BoolVarP(&fullNames,
		"full", "f", false, "show full names with namespace")
	Hostctl.AddCommand(listCmd)
}

var listCmd = &Command{
	Use:   "ls [pattern]",
	Short: "List hosts",
	Run: func(ctx *Context) {
		args := ctx.Args
		pattern := "*"
		if len(args) > 0 {
			pattern = args[0]
		}
		provider, err := providers.Get(providerName, true)
		fatal(err)
		for _, host := range provider.List(namespace + pattern) {
			if fullNames {
				fmt.Println(host.Name)
			} else {
				fmt.Println(strings.TrimPrefix(host.Name, namespace))
			}
		}
	},
}

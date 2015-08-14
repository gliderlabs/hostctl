package main

import (
	"fmt"
	"strings"

	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

var fullNames bool

func init() {
	listCmd.Flags().BoolVarP(&fullNames,
		"full", "f", false, "show full names with namespace")
	Hostctl.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "ls [pattern]",
	Short: "List hosts",
	Run: func(cmd *cobra.Command, args []string) {
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

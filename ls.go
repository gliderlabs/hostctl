package main

import (
	"fmt"

	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

func init() {
	HostctlCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "ls [pattern]",
	Short: "List hosts or hosts matching pattern",
	Run: func(cmd *cobra.Command, args []string) {
		pattern := "*"
		if len(args) > 0 {
			pattern = args[0]
		}
		provider, err := providers.Get(providerName)
		fatal(err)
		for _, host := range provider.List(namespace + pattern) {
			fmt.Println(host.Name)
		}
	},
}

package main

import (
	"fmt"
	"os"

	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

func init() {
	Hostctl.AddCommand(ipCmd)
}

var ipCmd = &cobra.Command{
	Use:   "ip <name>",
	Short: "Show IP for host",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 && defaultName == "" {
			cmd.Usage()
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

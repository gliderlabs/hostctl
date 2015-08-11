package main

import (
	"fmt"
	"os"

	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

func init() {
	HostctlCmd.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down <name>",
	Short: "Terminate a host if it exists",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 && defaultName == "" {
			cmd.Usage()
			os.Exit(1)
		}
		provider, err := providers.Get(providerName, true)
		fatal(err)
		name := fmt.Sprintf("%s%s", namespace, optArg(args, 0, defaultName))
		if provider.Get(name) == nil {
			return
		}
		fatal(provider.Destroy(name))
	},
}

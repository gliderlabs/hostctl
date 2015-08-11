package main

import (
	"fmt"
	"os"

	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

func init() {
	HostctlCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up <name>",
	Short: "Provision a host if it doesn't already exist",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			os.Exit(1)
		}
		provider, err := providers.Get(providerName, true)
		fatal(err)
		name := fmt.Sprintf("%s%s", namespace, args[0])
		if provider.Get(name) != nil {
			return
		}
		fatal(provider.Create(providers.Host{
			Name:    name,
			Flavor:  hostFlavor,
			Image:   hostImage,
			Region:  hostRegion,
			Keyname: hostKeyname,
		}))
	},
}

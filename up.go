package main

import (
	"os"
	"sync"

	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

func init() {
	HostctlCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up <name> [<name>...]",
	Short: "Provision a host if it doesn't already exist",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 && defaultName == "" {
			cmd.Usage()
			os.Exit(1)
		}
		if len(args) == 0 {
			args = []string{defaultName}
		}
		provider, err := providers.Get(providerName, true)
		fatal(err)
		parallelWait(args, func(_ int, arg string, wg *sync.WaitGroup) {
			defer wg.Done()
			name := namespace + arg
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
		})
	},
}

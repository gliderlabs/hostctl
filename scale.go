package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

func init() {
	HostctlCmd.AddCommand(scaleCmd)
}

var scaleCmd = &cobra.Command{
	Use:   "scale <name> <count>",
	Short: "Ensure a certain number of hosts are provisioned",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: support default name
		// TODO: support resizing a cluster
		// TODO: support scaling from a single host made with `up`
		// TODO: support scaling down
		if len(args) < 2 {
			cmd.Usage()
			os.Exit(1)
		}
		count, err := strconv.Atoi(args[1])
		fatal(err)
		var nodes []string
		for i := 0; i < count; i++ {
			nodes = append(nodes, fmt.Sprintf("%s%s.%v", namespace, args[0], i))
		}
		provider, err := providers.Get(providerName, true)
		fatal(err)
		parallelWait(nodes, func(_ int, node string, wg *sync.WaitGroup) {
			defer wg.Done()
			if provider.Get(node) != nil {
				return
			}
			fatal(provider.Create(providers.Host{
				Name:    node,
				Flavor:  hostFlavor,
				Image:   hostImage,
				Region:  hostRegion,
				Keyname: hostKeyname,
			}))
		})
	},
}

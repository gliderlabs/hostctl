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
		provider, err := providers.Get(providerName, true)
		fatal(err)
		nameFmt := namespace + args[0] + ".%v"
		var wg sync.WaitGroup
		for i := 0; i < count; i++ {
			wg.Add(1)
			name := fmt.Sprintf(nameFmt, i)
			go func() {
				defer wg.Done()
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
			}()
		}
		wg.Wait()
	},
}

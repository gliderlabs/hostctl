package main

import (
	"os"
	"sync"

	"github.com/gliderlabs/hostctl/providers"
	"github.com/spf13/cobra"
)

func init() {
	Hostctl.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up <name> [<name>...]",
	Short: "Provision host, wait until ready",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 && defaultName == "" {
			cmd.Usage()
			os.Exit(1)
		}
		loadStdinUserdata()
		if len(args) == 0 {
			args = []string{defaultName}
		}
		provider, err := providers.Get(providerName, true)
		fatal(err)
		finished := progressBar(".", 2)
		parallelWait(args, func(_ int, arg string, wg *sync.WaitGroup) {
			defer wg.Done()
			name := namespace + arg
			if provider.Get(name) != nil {
				return
			}
			fatal(provider.Create(newHost(name)))
		})
		finished <- true
	},
}

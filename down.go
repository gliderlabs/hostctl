package main

import (
	"os"
	"sync"

	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

func init() {
	Hostctl.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down <name> [<name>...]",
	Short: "Terminate host",
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
			if provider.Get(name) == nil {
				return
			}
			fatal(provider.Destroy(name))
		})
	},
}

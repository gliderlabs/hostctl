package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

func init() {
	HostctlCmd.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down <name> [<name>...]",
	Short: "Terminate a host if it exists",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 && defaultName == "" {
			cmd.Usage()
			os.Exit(1)
		}
		provider, err := providers.Get(providerName, true)
		fatal(err)
		count := len(args)
		if defaultName != "" && count == 0 {
			count = 1
		}
		var wg sync.WaitGroup
		for i := 0; i < count; i++ {
			wg.Add(1)
			name := fmt.Sprintf("%s%s", namespace, optArg(args, i, defaultName))
			go func() {
				defer wg.Done()
				if provider.Get(name) == nil {
					return
				}
				fatal(provider.Destroy(name))
			}()
		}
		wg.Wait()
	},
}

package main

import (
	"os"
	"sync"

	"github.com/gliderlabs/hostctl/providers"
)

func init() {
	Hostctl.AddCommand(downCmd)
}

var downCmd = &Command{
	Use:   "down <name> [<name>...]",
	Short: "Terminate host",
	Run: func(ctx *Context) {
		args := ctx.Args
		if len(args) < 1 && defaultName == "" {
			ctx.Cmd.Usage()
			os.Exit(1)
		}
		if len(args) == 0 {
			args = []string{defaultName}
		}
		provider, err := providers.Get(providerName, true)
		fatal(err)
		finished := progressBar(".", 1)
		parallelWait(args, func(_ int, arg string, wg *sync.WaitGroup) {
			defer wg.Done()
			name := namespace + arg
			if provider.Get(name) == nil {
				return
			}
			fatal(provider.Destroy(name))
		})
		finished <- true
	},
}

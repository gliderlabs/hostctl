package main

import (
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
		if len(ctx.Args) < 1 && defaultName == "" {
			ctx.Cmd.Usage()
			ctx.Exit(1)
		}
		if len(ctx.Args) == 0 {
			ctx.Args = []string{defaultName}
		}
		provider, err := providers.Get(providerName, true)
		ctxFatal(ctx, err)
		finished := progressBar(ctx, ".", 1)
		parallelWait(ctx.Args, func(_ int, arg string, wg *sync.WaitGroup) {
			defer wg.Done()
			name := namespace + arg
			if provider.Get(name) == nil {
				return
			}
			ctxFatal(ctx, provider.Destroy(name))
		})
		finished()
	},
}

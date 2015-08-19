package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/gliderlabs/hostctl/providers"
	"golang.org/x/crypto/ssh/terminal"
)

func newHost(name string) providers.Host {
	return providers.Host{
		Name:     name,
		Flavor:   hostFlavor,
		Image:    hostImage,
		Region:   hostRegion,
		Keyname:  hostKeyname,
		Userdata: hostUserdata,
	}
}

func loadStdinUserdata() {
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		data, err := ioutil.ReadAll(os.Stdin)
		fatal(err)
		hostUserdata = string(data)
	}
}

func parallelWait(items []string, fn func(int, string, *sync.WaitGroup)) {
	var wg sync.WaitGroup
	for i := 0; i < len(items); i++ {
		wg.Add(1)
		go fn(i, items[i], &wg)
	}
	wg.Wait()
}

func fatal(err error) {
	if err != nil {
		fmt.Println("!!", err)
		os.Exit(1)
	}
}

func ctxFatal(ctx *Context, err error) {
	if err != nil {
		fmt.Fprintln(ctx, "!!", err)
		ctx.Exit(1)
	}
}

func optArg(args []string, i int, default_ string) string {
	if i+1 > len(args) {
		return default_
	}
	return args[i]
}

func progressBar(ctx *Context, unit string, interval time.Duration) func() {
	finished := make(chan bool)
	go func() {
		for {
			select {
			case <-finished:
				return
			case <-time.After(interval * time.Second):
				fmt.Fprint(ctx.Err, unit)
			}
		}
	}()
	return func() {
		finished <- true
		fmt.Fprintln(ctx.Err)
	}
}

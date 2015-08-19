package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/gliderlabs/hostctl/providers"
)

func init() {
	Hostctl.AddCommand(scaleCmd)
}

var scaleCmd = &Command{
	Use:   "scale <name> <count>",
	Short: "Resize host cluster",
	Run: func(ctx *Context) {
		args := ctx.Args
		if (len(args) < 2 && defaultName == "") ||
			(len(args) < 1 && defaultName != "") {
			ctx.Cmd.Usage()
			os.Exit(1)
		}
		var name, count string
		if len(args) == 1 {
			name = defaultName
			count = args[0]
		} else {
			name = args[0]
			count = args[1]
		}
		loadStdinUserdata()
		provider, err := providers.Get(providerName, true)
		fatal(err)
		existing := existingHosts(provider, name)
		desired := desiredHosts(name, count)
		hosts := append(strSet(existing, desired), namespace+name)
		finished := progressBar(ctx, ".", 2)
		parallelWait(hosts, func(_ int, host string, wg *sync.WaitGroup) {
			defer wg.Done()
			if !strIn(host, desired) {
				fatal(provider.Destroy(host))
				return
			}
			if strIn(host, desired) && !strIn(host, existing) {
				fatal(provider.Create(newHost(host)))
				return
			}
		})
		finished()
	},
}

func desiredHosts(name, count string) []string {
	c, err := strconv.Atoi(count)
	fatal(err)
	var hosts []string
	for i := 0; i < c; i++ {
		hosts = append(hosts, fmt.Sprintf("%s%s.%v", namespace, name, i))
	}
	return hosts
}

func existingHosts(provider providers.HostProvider, name string) []string {
	var hosts []string
	for _, h := range provider.List(namespace + name + ".*") {
		hosts = append(hosts, h.Name)
	}
	return hosts
}

func strIn(str string, list []string) bool {
	for i := range list {
		if str == list[i] {
			return true
		}
	}
	return false
}

func strSet(strs ...[]string) []string {
	m := make(map[string]bool)
	for i := range strs {
		for _, str := range strs[i] {
			m[str] = true
		}
	}
	var set []string
	for k := range m {
		set = append(set, k)
	}
	return set
}

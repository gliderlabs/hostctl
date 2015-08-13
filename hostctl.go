package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/MattAitchison/env"
	"github.com/spf13/cobra"

	_ "github.com/progrium/hostctl/digitalocean"
)

var (
	providerName = env.String("HOSTCTL_PROVIDER", "digitalocean", "cloud provider")
	defaultName  = env.String("HOSTCTL_NAME", "", "optional default name")
	namespace    = env.String("HOSTCTL_NAMESPACE", "", "optional namespace for names")

	hostImage   = env.String("HOSTCTL_IMAGE", "", "vm image")
	hostFlavor  = env.String("HOSTCTL_FLAVOR", "", "vm flavor")
	hostRegion  = env.String("HOSTCTL_REGION", "", "vm region")
	hostKeyname = env.String("HOSTCTL_KEYNAME", "", "vm keyname")

	user = env.String("HOSTCTL_USER", os.Getenv("USER"), "ssh user")
)

var HostctlCmd = &cobra.Command{
	Use:   "hostctl",
	Short: "An opinionated tool for provisioning cloud VMs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
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

func optArg(args []string, i int, default_ string) string {
	if i+1 > len(args) {
		return default_
	}
	return args[i]
}

func main() {
	HostctlCmd.Execute()
}

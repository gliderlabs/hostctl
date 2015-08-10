package main

import (
	"fmt"
	"os"

	"github.com/MattAitchison/envconfig"
	"github.com/spf13/cobra"

	_ "github.com/progrium/hostctl/digitalocean"
	"github.com/progrium/hostctl/providers"
)

var (
	providerName = envconfig.String("HOSTCTL_PROVIDER", "digitalocean", "cloud provider to use")
	defaultName  = envconfig.String("HOSTCTL_NAME", "", "default name to use")
	namespace    = envconfig.String("HOSTCTL_NAMESPACE", "", "namespace to use")

	hostImage   = envconfig.String("HOSTCTL_IMAGE", "", "image to use for vm")
	hostFlavor  = envconfig.String("HOSTCTL_FLAVOR", "", "flavor to use for vm")
	hostRegion  = envconfig.String("HOSTCTL_REGION", "", "region to use for vm")
	hostKeyname = envconfig.String("HOSTCTL_KEYNAME", "", "comma-separated list of keynames")
)

var HostctlCmd = &cobra.Command{
	Use:   "hostctl",
	Short: "An opinionated tool for provisioning cloud VMs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func fatal(err error) {
	if err != nil {
		fmt.Println("!!", err)
		os.Exit(1)
	}
}

func hostExists(provider providers.HostProvider, name string) bool {
	hosts := provider.List(namespace + "*")
	for i := range hosts {
		if hosts[i].Name == name {
			return true
		}
	}
	return false
}

func main() {
	HostctlCmd.Execute()
}

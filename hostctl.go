package main

import (
	"fmt"
	"os"

	"github.com/progrium/envconfig"
	"github.com/spf13/cobra"

	_ "github.com/progrium/hostctl/digitalocean"
)

var (
	providerName = envconfig.String("HOSTCTL_PROVIDER", "digitalocean", "cloud provider")
	defaultName  = envconfig.String("HOSTCTL_NAME", "", "optional default name")
	namespace    = envconfig.String("HOSTCTL_NAMESPACE", "", "optional namespace for names")

	hostImage   = envconfig.String("HOSTCTL_IMAGE", "", "vm image")
	hostFlavor  = envconfig.String("HOSTCTL_FLAVOR", "", "vm flavor")
	hostRegion  = envconfig.String("HOSTCTL_REGION", "", "vm region")
	hostKeyname = envconfig.String("HOSTCTL_KEYNAME", "", "vm keyname")
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

func main() {
	HostctlCmd.Execute()
}

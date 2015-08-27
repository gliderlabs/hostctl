package main

import (
	"fmt"
	"os"

	"github.com/MattAitchison/env"
	"github.com/spf13/cobra"
)

var (
	Version string

	providerName string
	defaultName  string
	namespace    string

	hostImage    string
	hostFlavor   string
	hostRegion   string
	hostKeyname  string
	hostUserdata string

	sshUser string

	profile string
)

func readEnv() {
	env.Clear()
	providerName = env.String("HOSTCTL_PROVIDER", "digitalocean", "cloud provider")
	defaultName = env.String("HOSTCTL_NAME", "", "optional default name")
	namespace = env.String("HOSTCTL_NAMESPACE", "", "optional namespace for names")
	hostImage = env.String("HOSTCTL_IMAGE", "", "vm image")
	hostFlavor = env.String("HOSTCTL_FLAVOR", "", "vm flavor")
	hostRegion = env.String("HOSTCTL_REGION", "", "vm region")
	hostKeyname = env.String("HOSTCTL_KEYNAME", "", "vm keyname")
	hostUserdata = env.String("HOSTCTL_USERDATA", "", "optional vm user data")
	sshUser = env.String("HOSTCTL_USER", os.Getenv("USER"), "ssh user")
}

func init() {
	readEnv()
	Hostctl.PersistentFlags().StringVarP(&profile, "profile", "p", "", "profile to source")
	Hostctl.AddCommand(versionCmd)
}

func main() {
	fatal(Hostctl.Execute())
}

var Hostctl = &cobra.Command{
	Use:   "hostctl",
	Short: "An opinionated tool for provisioning cloud VMs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if exists(expandHome("~/.hostctl")) {
			fatal(source(expandHome("~/.hostctl")))
			readEnv()
		}
		if profile != "" {
			fatal(source(profile))
			readEnv()
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

package main

import (
	"os"

	"github.com/progrium/envconfig"
	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

var (
	exportMode  bool
	secretsMode bool
)

func init() {
	envCmd.Flags().BoolVarP(&exportMode,
		"export", "e", false, "export vars for sourcing later")
	envCmd.Flags().BoolVarP(&secretsMode,
		"secrets", "s", false, "show secrets or include in export")
	HostctlCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Shows current relevant environment",
	Run: func(cmd *cobra.Command, args []string) {
		envconfig.PrintEnv(os.Stdout, exportMode, secretsMode)
		provider, _ := providers.Get(providerName, false)
		if provider != nil {
			provider.Env().PrintEnv(os.Stdout, exportMode, secretsMode)
		}
	},
}

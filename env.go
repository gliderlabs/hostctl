package main

import (
	"os"

	"github.com/MattAitchison/env"
	"github.com/gliderlabs/hostctl/providers"
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
	Hostctl.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show relevant environment",
	Run: func(cmd *cobra.Command, args []string) {
		env.PrintEnv(os.Stdout, exportMode, secretsMode)
		provider, _ := providers.Get(providerName, false)
		if provider != nil && provider.Env() != nil {
			provider.Env().PrintEnv(os.Stdout, exportMode, secretsMode)
		}
	},
}

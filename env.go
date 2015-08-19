package main

import (
	"github.com/MattAitchison/env"
	"github.com/gliderlabs/hostctl/providers"
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

var envCmd = &Command{
	Use:   "env",
	Short: "Show relevant environment",
	Run: func(ctx *Context) {
		env.PrintEnv(ctx, exportMode, secretsMode)
		provider, _ := providers.Get(providerName, false)
		if provider != nil && provider.Env() != nil {
			provider.Env().PrintEnv(ctx, exportMode, secretsMode)
		}
	},
}

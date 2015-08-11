package main

import (
	"fmt"
	"io"

	"github.com/progrium/envconfig"
	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

var (
	exportMode  bool
	secretsMode bool
)

func init() {
	envCmd.Flags().BoolVarP(&exportMode, "export", "e", false, "export vars for sourcing later")
	envCmd.Flags().BoolVarP(&secretsMode, "secrets", "s", false, "show secrets or include in export")
	HostctlCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Shows current relevant environment",
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range envconfig.Vars() {
			printVar(v)
		}
		provider, _ := providers.Get(providerName, false)
		if provider != nil {
			for _, v := range provider.Env().Vars() {
				printVar(v)
			}
		}
	},
}

func printVar(out io.Writer, v *envconfig.ConfigVar) {
	value := v.Value.String()
	if v.Secret {
		if secretsMode {
			value = v.Value.Get().(string)
		} else {
			if exportMode {
				return
			}
		}
	}
	if exportMode {
		fmt.Fprintf(out, "export %s=\"%s\"\n", v.Name, value)
	} else {
		kv := fmt.Sprintf("%s=\"%s\"", v.Name, value)
		fmt.Fprintf(out, "%-40s # %s\n", kv, v.Description)
	}
}

package main

import (
	"fmt"

	"github.com/MattAitchison/envconfig"
	"github.com/spf13/cobra"
)

func init() {
	HostctlCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Shows current configuration / environment",
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range envconfig.Vars() {
			fmt.Printf("export %s=\"%s\" # %s\n", v.Name, v.Value.Get(), v.Description)
		}
	},
}

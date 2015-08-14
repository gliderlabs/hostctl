package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/progrium/hostctl/providers"
	"github.com/spf13/cobra"
)

func init() {
	Hostctl.AddCommand(sshCmd)
}

var sshCmd = &cobra.Command{
	Use:   "ssh <name> [--] [<command>...]",
	Short: "SSH to host",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 && defaultName == "" {
			cmd.Usage()
			os.Exit(1)
		}
		name, sshCmd := sshParseArgs(args)
		provider, err := providers.Get(providerName, true)
		fatal(err)
		host := provider.Get(namespace + name)
		if host == nil {
			os.Exit(1)
		}
		fatal(sshExec(host.IP, sshCmd))
	},
}

func sshExec(ip string, cmd []string) error {
	binary, err := exec.LookPath("ssh")
	if err != nil {
		return fmt.Errorf("Unable to find ssh")
	}
	args := []string{"ssh", "-A", "-l", user, ip}
	return syscall.Exec(binary, append(args, cmd...), os.Environ())
}

func sshParseArgs(args []string) (string, []string) {
	var name string
	var sshCmd []string
	if len(args) == 0 || args[0] == "--" {
		name = defaultName
		sshCmd = args
	} else {
		name = args[0]
		sshCmd = args[1:]
	}
	if optArg(sshCmd, 0, "") == "--" {
		sshCmd = sshCmd[1:]
	}
	return name, sshCmd
}

package main

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Command struct {
	Use     string
	Aliases []string
	Short   string
	Long    string
	Example string
	Run     func(*Context)

	Cmd     *cobra.Command
	Context Context
}

func (c *Command) Execute() error {
	return c.setup().Cmd.Execute()
}

func (c *Command) AddCommand(cmd *Command) {
	cmd.Context = c.setup().Context
	c.Cmd.AddCommand(cmd.setup().Cmd)
}

func (c *Command) Flags() *pflag.FlagSet {
	return c.setup().Cmd.Flags()
}

func (c *Command) setup() *Command {
	if c.Cmd != nil {
		return c
	}
	cmd := &cobra.Command{
		Use:     c.Use,
		Aliases: c.Aliases,
		Short:   c.Short,
		Long:    c.Long,
		Example: c.Example,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := c.Context // make a copy
			ctx.Args = args
			ctx.Cmd = cmd
			cmd.SetOutput(&ctx)
			c.Run(&ctx)
		},
	}
	c.Cmd = cmd
	return c
}

type Context struct {
	Out  io.Writer
	Err  io.Writer
	In   io.Reader
	Exit func(int)
	Args []string
	Cmd  *cobra.Command
}

func (c *Context) Read(data []byte) (n int, err error) {
	return c.In.Read(data)
}

func (c *Context) Write(data []byte) (n int, err error) {
	return c.Out.Write(data)
}

func (c *Context) Arg(i int) string {
	return c.OptArg(i, "")
}

func (c *Context) OptArg(i int, default_ string) string {
	if i+1 > len(c.Args) {
		return default_
	}
	return c.Args[i]
}

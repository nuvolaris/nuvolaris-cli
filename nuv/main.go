package main

import (
	"github.com/alecthomas/kong"
)

type CLI struct {
	Task TaskCmd `cmd:"" help:"task subcommand."`
	Wsk  WskCmd  `cmd:"" help:"wsk subcommand."`
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("nuv"),
		kong.Description("nuv is the command line tool to manage Nuvolaris"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: true,
		}),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

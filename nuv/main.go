package main

import (
	"github.com/alecthomas/kong"
)

// CLI_VERSION holds the current version, to be set by the build with
//  go build -ldflags "-X main.CLI_VERSION=<version>"
var CLI_VERSION string = "latest"

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

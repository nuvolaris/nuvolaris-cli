package main

type CLI struct {
	WskCLI
	TaskCLI
}

type WskCLI struct {
	Wsk struct {
	} `cmd:"" help:"wsk subcommand."`
}

type TaskCLI struct {
	Task struct {
	} `cmd:"" help:"task subcommand."`
}

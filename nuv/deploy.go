package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
)

//go:embed embed/nuvolaris.yml
var NuvolarisYml []byte

type DeployCmd struct {
	Args []string `optional:"" name:"args" help:"kind subcommand args"`
}

func (*DeployCmd) Run() error {
	fmt.Println("Deploying Nuvolaris...")
	ioutil.WriteFile("nuvolaris.yml", NuvolarisYml, 0600)
	Task()
	return nil
}

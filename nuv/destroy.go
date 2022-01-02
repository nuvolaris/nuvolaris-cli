package main

type DestroyCmd struct {
	Args []string `optional:"" name:"args" help:"destroy nuvolaris cluster"`
}

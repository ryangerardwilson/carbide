package main

import "github.com/ryangerardwilson/carbide/cli/internal/cli"

var commit string

func main() {
	cli.SetCommit(commit)
	cli.Main()
}

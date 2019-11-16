package main

import "github.com/arkste/coolctl/cmd"

// Version represents the app version
var Version string

func main() {
	cmd.Execute(Version)
}

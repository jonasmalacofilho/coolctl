package main

import "github.com/arkste/coolctl/nzxt"

func main() {
	kraken := nzxt.NewDriver()
	kraken.Connect()
	defer kraken.Disconnect()

	kraken.GetStatus()
}

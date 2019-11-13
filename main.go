package main

import "github.com/arkste/coolctl/driver"

func main() {
	kraken := driver.NewKrakenDriver()
	kraken.Connect()
	defer kraken.Disconnect()

	kraken.GetStatus()
}

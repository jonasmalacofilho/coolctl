package main

import "github.com/arkste/coolctl/driver"

func main() {
	kraken := driver.NewKrakenDriver()
	kraken.Connect()
	kraken.GetStatus()
	kraken.SetColor("ring", "loading", []string{"0000FF"})
	kraken.SetColor("logo", "fixed", []string{"0000FF"})
}

package main

import "github.com/arkste/coolctl/driver"

func main() {
	kraken := driver.NewKrakenDriver()
	kraken.Connect()
	kraken.GetStatus()

	// kraken.SetColor("ring", "loading", []string{"FF0000"})
	// kraken.SetColor("logo", "pulse", []string{"FF0000"})
	// kraken.SetColor("sync", "super-fixed", []string{"FF0000", "00FF00", "0000FF"})

	kraken.SetColor("logo", "off", nil)
	kraken.SetColor("ring", "fading", []string{"FF0000", "00FF00", "0000FF"})
	kraken.SetSpeed("pump", "20 60  35 60  55 100  60 100")
	kraken.SetSpeed("fan", "20 25  35 25  50 55  60 100")
}

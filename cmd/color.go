package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arkste/coolctl/driver"
)

// colorCmd represents the color command
var colorCmd = &cobra.Command{
	Use:   "color",
	Short: "set the color of the logo, ring or sync",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a color channel (e.g: logo, ring or sync)")
		}

		if len(args) < 2 {
			return errors.New("requires a color mode (e.g: off, fading, etc)")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		kraken := driver.NewKrakenDriver()
		kraken.Connect()
		kraken.SetColor(args[0], args[1], args[2:])

		fmt.Printf("%s: %s %s", args[0], args[1], args[2:])
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(colorCmd)
}

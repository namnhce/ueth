package main

import (
	"os"

	"github.com/gookit/color"
	"github.com/urfave/cli"

	"github.com/namnhce/ueth/pkg/wallets"
)

func main() {
	app := cli.NewApp()
	app.Name = "ueth"
	app.Usage = "ETH wallet utilities tool"

	app.Commands = []cli.Command{
		{
			Name:  "wallet",
			Usage: "wallet seed commands",
			Subcommands: []cli.Command{
				{
					Name:  "gen",
					Usage: "generate wallet from mnemonic seed phrase",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "mnemonic, m",
							Usage:    "Mnemonic seed phrase",
							Required: false,
						},
						cli.StringFlag{
							Name:     "num-wallets, n",
							Usage:    "Numbers of wallets to generate",
							Required: false,
						},
						cli.StringFlag{
							Name:     "output, o",
							Usage:    "CSV Output Filename",
							Required: false,
						},
					},
					Action: func(c *cli.Context) error {
						return displayErrorOrMessage(wallets.DoGenerateWallet(c))
					},
				},
				{
					Name:  "send",
					Usage: "send ETH to list of addresses in CSV file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "private-key, k",
							Usage:    "private key of sender",
							Required: false,
						},
						cli.StringFlag{
							Name:     "value, v",
							Usage:    "the ETH value to send",
							Required: false,
						},
						cli.StringFlag{
							Name:     "input, i",
							Usage:    "CSV file path",
							Required: false,
						},
					},
					Action: func(c *cli.Context) error {
						return displayErrorOrMessage(wallets.DoSend(c))
					},
				},
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		app.Command("help").Run(c)
		return nil
	}

	app.Run(os.Args)
}

func displayErrorOrMessage(err error) error {
	if err != nil {
		color.Error.Tips(err.Error())
		return cli.NewExitError(color.Error.Sprintf(err.Error()), 1)
	}

	return nil
}

package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/bryanl/doit"
	"github.com/codegangsta/cli"
)

func accountCommands() cli.Command {
	return cli.Command{
		Name:  "account",
		Usage: "account commands",
		Action: func(c *cli.Context) {
			config := doit.NewCLIConfig(c.GlobalString("token"), c.App.Writer)
			a, err := doit.AccountGet(config)
			if err != nil {
				logrus.WithField("err", err).Error("could not display account")
				return
			}

			doit.WriteJSON(a, config.Writer())
		},
	}
}

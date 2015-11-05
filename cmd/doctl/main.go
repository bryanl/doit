package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/digitalocean/doctl/commands"
)

func init() {
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	commands.Execute()
}

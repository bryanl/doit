package commands

import (
	"io"

	"github.com/digitalocean/doctl"
	"github.com/spf13/cobra"
)

// Account creates the account commands heirarchy.
func Account() *cobra.Command {
	cmdAccount := &cobra.Command{
		Use:   "account",
		Short: "account commands",
		Long:  "account is used to access account commands",
	}

	cmdAccountGet := cmdBuilder(RunAccountGet, "get", "get account", writer, "g")
	cmdAccount.AddCommand(cmdAccountGet)

	return cmdAccount
}

// RunAccountGet runs account get.
func RunAccountGet(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()

	a, _, err := client.Account.Get()
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(a, out)
}

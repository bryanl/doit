package commands

import (
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

// ImageAction creates the image action commmand.
func ImageAction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image-action",
		Short: "image-action commands",
		Long:  "image-action commands",
	}

	cmdImageActionsGet := cmdBuilder(RunImageActionsGet,
		"get", "get image action", writer)
	cmd.AddCommand(cmdImageActionsGet)
	addIntFlag(cmdImageActionsGet, doctl.ArgImageID, 0, "image id")
	addIntFlag(cmdImageActionsGet, doctl.ArgActionID, 0, "action id")

	cmdImageActionsTransfer := cmdBuilder(RunImageActionsTransfer,
		"transfer", "transfer imagr", writer)
	cmd.AddCommand(cmdImageActionsTransfer)
	addIntFlag(cmdImageActionsTransfer, doctl.ArgImageID, 0, "image id")
	addStringFlag(cmdImageActionsTransfer, doctl.ArgRegionSlug, "", "region")

	return cmd
}

// RunImageActionsGet retrieves an action for an image.
func RunImageActionsGet(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	imageID := doctl.DoctlConfig.GetInt(ns, doctl.ArgImageID)
	actionID := doctl.DoctlConfig.GetInt(ns, doctl.ArgActionID)

	action, _, err := client.ImageActions.Get(imageID, actionID)
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(action, out)
}

// RunImageActionsTransfer an image.
func RunImageActionsTransfer(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetInt(ns, doctl.ArgImageID)
	req := &godo.ActionRequest{
		"region": doctl.DoctlConfig.GetString(ns, doctl.ArgRegionSlug),
	}

	action, _, err := client.ImageActions.Transfer(id, req)
	if err != nil {
		logrus.WithField("err", err).Fatal("could not transfer image")
	}

	return doctl.DisplayOutput(action, out)
}

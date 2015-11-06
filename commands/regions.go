package commands

import (
	"io"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

// Region creates the region commands heirarchy.
func Region() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "region",
		Short: "region commands",
		Long:  "region is used to access region commands",
	}

	cmdRegionList := cmdBuilder(RunRegionList, "list", "list regions", writer)
	cmd.AddCommand(cmdRegionList)

	return cmd
}

// RunRegionList all regions.
func RunRegionList(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()

	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.Regions.List(opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := doctl.PaginateResp(f)
	if err != nil {
		return err
	}

	list := make([]godo.Region, len(si))
	for i := range si {
		list[i] = si[i].(godo.Region)
	}

	return doctl.DisplayOutput(list, out)
}

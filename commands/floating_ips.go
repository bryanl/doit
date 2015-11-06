package commands

import (
	"errors"
	"io"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

// FloatingIP creates the command heirarchy for floating ips.
func FloatingIP() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "floating-ip",
		Short:   "floating IP commands",
		Long:    "floating-ip is used to access commands on floating IPs",
		Aliases: []string{"fip"},
	}

	cmdFloatingIPCreate := cmdBuilder(RunFloatingIPCreate, "create", "create a floating IP", writer, "c")
	cmd.AddCommand(cmdFloatingIPCreate)
	addStringFlag(cmdFloatingIPCreate, doctl.ArgRegionSlug, "", "Region where to create the floating IP.")
	addIntFlag(cmdFloatingIPCreate, doctl.ArgDropletID, 0, "ID of the droplet to assign the IP to. (Optional)")

	cmdFloatingIPGet := cmdBuilder(RunFloatingIPGet, "get", "get the details of a floating IP", writer, "g")
	cmd.AddCommand(cmdFloatingIPGet)
	addStringFlag(cmdFloatingIPGet, doctl.ArgIPAddress, "", "IP address of the floating IP")

	cmdFloatingIPDelete := cmdBuilder(RunFloatingIPDelete, "delete", "delete a floating IP address", writer, "d")
	cmd.AddCommand(cmdFloatingIPDelete)
	addStringFlag(cmdFloatingIPDelete, doctl.ArgIPAddress, "", "IP address of the floating IP")

	cmdFloatingIPList := cmdBuilder(RunFloatingIPList, "list", "list all floating IP addresses", writer, "ls")
	cmd.AddCommand(cmdFloatingIPList)

	return cmd
}

// RunFloatingIPCreate runs floating IP create.
func RunFloatingIPCreate(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	req := &godo.FloatingIPCreateRequest{
		Region:    doctl.DoctlConfig.GetString(ns, doctl.ArgRegionSlug),
		DropletID: doctl.DoctlConfig.GetInt(ns, doctl.ArgDropletID),
	}
	ip, _, err := client.FloatingIPs.Create(req)
	if err != nil {
		return err
	}
	return doctl.DisplayOutput(ip, out)
}

// RunFloatingIPGet retrieves a floating IP's details.
func RunFloatingIPGet(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	ip := doctl.DoctlConfig.GetString(ns, doctl.ArgIPAddress)

	if len(ip) < 1 {
		return errors.New("invalid ip address")
	}

	d, _, err := client.FloatingIPs.Get(ip)
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(d, out)
}

// RunFloatingIPDelete runs floating IP delete.
func RunFloatingIPDelete(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	ip := doctl.DoctlConfig.GetString(ns, doctl.ArgIPAddress)
	_, err := client.FloatingIPs.Delete(ip)
	return err
}

// RunFloatingIPList runs floating IP create.
func RunFloatingIPList(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()

	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.FloatingIPs.List(opt)
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

	list := make([]godo.FloatingIP, len(si))
	for i := range si {
		list[i] = si[i].(godo.FloatingIP)
	}

	return doctl.DisplayOutput(list, out)
}

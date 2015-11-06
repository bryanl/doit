package commands

import (
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/digitalocean/godo/util"
	"github.com/spf13/cobra"
)

// Droplet creates the droplet command.
func Droplet() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "droplet",
		Aliases: []string{"d"},
		Short:   "droplet commands",
		Long:    "droplet is used to access droplet commands",
	}

	cmdDropletActions := cmdBuilder(RunDropletActions,
		"actions", "droplet actions", writer, "a")
	cmd.AddCommand(cmdDropletActions)
	addIntFlag(cmdDropletActions, doctl.ArgDropletID, 0, "Droplet ID")

	cmdDropletBackups := cmdBuilder(RunDropletBackups,
		"backups", "droplet backups", writer, "b")
	cmd.AddCommand(cmdDropletBackups)
	addIntFlag(cmdDropletBackups, doctl.ArgDropletID, 0, "Droplet ID")

	cmdDropletCreate := cmdBuilder(RunDropletCreate,
		"create", "create droplet", writer, "c")
	cmd.AddCommand(cmdDropletCreate)
	addStringSliceFlag(cmdDropletCreate, doctl.ArgSSHKeys, []string{}, "SSH Keys or fingerprints")
	addStringFlag(cmdDropletCreate, doctl.ArgUserData, "", "User data")
	addStringFlag(cmdDropletCreate, doctl.ArgUserDataFile, "", "User data file")
	addBoolFlag(cmdDropletCreate, doctl.ArgDropletWait, false, "Wait for droplet to be created")
	addStringFlag(cmdDropletCreate, doctl.ArgDropletName, "", "Droplet name")
	addStringFlag(cmdDropletCreate, doctl.ArgRegionSlug, "", "Droplet region")
	addStringFlag(cmdDropletCreate, doctl.ArgSizeSlug, "", "Droplet size")
	addBoolFlag(cmdDropletCreate, doctl.ArgBackups, false, "Backup droplet")
	addBoolFlag(cmdDropletCreate, doctl.ArgIPv6, false, "IPv6 support")
	addBoolFlag(cmdDropletCreate, doctl.ArgPrivateNetworking, false, "Private networking")
	addStringFlag(cmdDropletCreate, doctl.ArgImage, "", "Droplet image")

	cmdDropletDelete := cmdBuilder(RunDropletDelete,
		"delete", "delete droplet", writer, "d")
	cmd.AddCommand(cmdDropletDelete)
	addIntFlag(cmdDropletDelete, doctl.ArgDropletID, 0, "Droplet ID")

	cmdDropletGet := cmdBuilder(RunDropletGet,
		"get", "get droplet", writer, "g")
	cmd.AddCommand(cmdDropletGet)
	addIntFlag(cmdDropletGet, doctl.ArgDropletID, 0, "Droplet ID")

	cmdDropletKernels := cmdBuilder(RunDropletKernels,
		"kernels", "droplet kernels", writer, "k")
	cmd.AddCommand(cmdDropletKernels)
	addIntFlag(cmdDropletKernels, doctl.ArgDropletID, 0, "Droplet ID")

	cmdDropletList := cmdBuilder(RunDropletList,
		"list", "list droplets", writer, "ls")
	cmd.AddCommand(cmdDropletList)

	cmdDropletNeighbors := cmdBuilder(RunDropletNeighbors,
		"neighbors", "droplet neighbors", writer, "n")
	cmd.AddCommand(cmdDropletNeighbors)
	addIntFlag(cmdDropletNeighbors, doctl.ArgDropletID, 0, "Droplet ID")

	cmdDropletSnapshots := cmdBuilder(RunDropletSnapshots,
		"snapshots", "snapshots", writer, "s")
	cmd.AddCommand(cmdDropletSnapshots)
	addIntFlag(cmdDropletSnapshots, doctl.ArgDropletID, 0, "Droplet ID")

	return cmd
}

// NewCmdDropletActions creates a droplet action get command.
func NewCmdDropletActions(out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "actions",
		Short: "get droplet actions",
		Long:  "get droplet actions",
		Run: func(cmd *cobra.Command, args []string) {
			checkErr(RunDropletActions(cmdNS(cmd), out), cmd)
		},
	}
}

// RunDropletActions returns a list of actions for a droplet.
func RunDropletActions(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetInt(ns, doctl.ArgDropletID)

	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.Droplets.Actions(id, opt)
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

	list := make([]godo.Action, len(si))
	for i := range si {
		list[i] = si[i].(godo.Action)
	}

	return doctl.DisplayOutput(list, out)
}

// RunDropletBackups returns a list of backup images for a droplet.
func RunDropletBackups(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetInt(ns, doctl.ArgDropletID)

	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.Droplets.Backups(id, opt)
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

	list := make([]godo.Image, len(si))
	for i := range si {
		list[i] = si[i].(godo.Image)
	}

	return doctl.DisplayOutput(list, out)
}

// RunDropletCreate creates a droplet.
func RunDropletCreate(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()

	sshKeys := []godo.DropletCreateSSHKey{}
	for _, rawKey := range doctl.DoctlConfig.GetStringSlice(ns, doctl.ArgSSHKeys) {
		rawKey = strings.TrimPrefix(rawKey, "[")
		rawKey = strings.TrimSuffix(rawKey, "]")
		if i, err := strconv.Atoi(rawKey); err == nil {
			sshKeys = append(sshKeys, godo.DropletCreateSSHKey{ID: i})
			continue
		}

		sshKeys = append(sshKeys, godo.DropletCreateSSHKey{Fingerprint: rawKey})
	}

	userData := doctl.DoctlConfig.GetString(ns, doctl.ArgUserData)
	if userData == "" && doctl.DoctlConfig.GetString(ns, doctl.ArgUserDataFile) != "" {
		data, err := ioutil.ReadFile(doctl.DoctlConfig.GetString(ns, doctl.ArgUserDataFile))
		if err != nil {
			return err
		}
		userData = string(data)
	}

	wait := doctl.DoctlConfig.GetBool(ns, doctl.ArgDropletWait)

	dcr := &godo.DropletCreateRequest{
		Name:              doctl.DoctlConfig.GetString(ns, doctl.ArgDropletName),
		Region:            doctl.DoctlConfig.GetString(ns, doctl.ArgRegionSlug),
		Size:              doctl.DoctlConfig.GetString(ns, doctl.ArgSizeSlug),
		Backups:           doctl.DoctlConfig.GetBool(ns, doctl.ArgBackups),
		IPv6:              doctl.DoctlConfig.GetBool(ns, doctl.ArgIPv6),
		PrivateNetworking: doctl.DoctlConfig.GetBool(ns, doctl.ArgPrivateNetworking),
		SSHKeys:           sshKeys,
		UserData:          userData,
	}

	imageStr := doctl.DoctlConfig.GetString(ns, doctl.ArgImage)
	if i, err := strconv.Atoi(imageStr); err == nil {
		dcr.Image = godo.DropletCreateImage{ID: i}
	} else {
		dcr.Image = godo.DropletCreateImage{Slug: imageStr}
	}

	r, resp, err := client.Droplets.Create(dcr)
	if err != nil {
		return err
	}

	var action *godo.LinkAction

	if wait {
		for _, a := range resp.Links.Actions {
			if a.Rel == "create" {
				action = &a
			}
		}
	}

	if action != nil {
		err = util.WaitForActive(client, action.HREF)
		if err != nil {
			return err
		}

		r, err = getDropletByID(client, r.ID)
		if err != nil {
			return err
		}
	}

	return doctl.DisplayOutput(r, out)
}

// RunDropletDelete destroy a droplet by id.
func RunDropletDelete(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetInt(ns, doctl.ArgDropletID)

	_, err := client.Droplets.Delete(id)
	return err
}

// RunDropletGet returns a droplet.
func RunDropletGet(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetInt(ns, doctl.ArgDropletID)

	droplet, err := getDropletByID(client, id)
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(droplet, out)
}

// RunDropletKernels returns a list of available kernels for a droplet.
func RunDropletKernels(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetInt(ns, doctl.ArgDropletID)

	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.Droplets.Kernels(id, opt)
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

	list := make([]godo.Kernel, len(si))
	for i := range si {
		list[i] = si[i].(godo.Kernel)
	}

	return doctl.DisplayOutput(list, out)
}

// RunDropletList returns a list of droplets.
func RunDropletList(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()

	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.Droplets.List(opt)
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

	list := make([]godo.Droplet, len(si))
	for i := range si {
		list[i] = si[i].(godo.Droplet)
	}

	return doctl.DisplayOutput(list, out)
}

// RunDropletNeighbors returns a list of droplet neighbors.
func RunDropletNeighbors(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetInt(ns, doctl.ArgDropletID)

	list, _, err := client.Droplets.Neighbors(id)
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(list, out)
}

// RunDropletSnapshots returns a list of available kernels for a droplet.
func RunDropletSnapshots(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetInt(ns, doctl.ArgDropletID)

	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := client.Droplets.Snapshots(id, opt)
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

	list := make([]godo.Image, len(si))
	for i := range si {
		list[i] = si[i].(godo.Image)
	}

	return doctl.DisplayOutput(list, out)
}

func getDropletByID(client *godo.Client, id int) (*godo.Droplet, error) {
	if id < 1 {
		return nil, errors.New("missing droplet id")
	}

	droplet, _, err := client.Droplets.Get(id)
	return droplet, err
}

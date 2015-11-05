package commands

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
)

// Images creates an image command.
func Images() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "image commands",
		Long:  "image commands",
	}

	out := os.Stdout

	cmdImagesList := cmdBuilder(RunImagesList, "list", "list images", out)
	cmd.AddCommand(cmdImagesList)
	addBoolFlag(cmdImagesList, doctl.ArgImagePublic, false, "List public images")

	cmdImagesListDistribution := cmdBuilder(RunImagesListDistribution,
		"list-distribution", "list distribution images", out)
	cmd.AddCommand(cmdImagesListDistribution)
	addBoolFlag(cmdImagesListDistribution, doctl.ArgImagePublic, false, "List public images")

	cmdImagesListApplication := cmdBuilder(RunImagesListDistribution,
		"list-application", "list application images", out)
	cmd.AddCommand(cmdImagesListApplication)
	addBoolFlag(cmdImagesListApplication, doctl.ArgImagePublic, false, "List public images")

	cmdImagesListUser := cmdBuilder(RunImagesListDistribution,
		"list-user", "list user images", out)
	cmd.AddCommand(cmdImagesListUser)
	addBoolFlag(cmdImagesListUser, doctl.ArgImagePublic, false, "List public images")

	cmdImagesGet := cmdBuilder(RunImagesGet, "get", "Get image", out)
	cmd.AddCommand(cmdImagesGet)
	addStringFlag(cmdImagesGet, doctl.ArgImage, "", "Image id")

	cmdImagesUpdate := cmdBuilder(RunImagesUpdate, "update", "Update image", out)
	cmd.AddCommand(cmdImagesUpdate)
	addStringFlag(cmdImagesUpdate, doctl.ArgImage, "", "Image id")
	addStringFlag(cmdImagesUpdate, doctl.ArgImageName, "", "Image name")

	cmdImagesDelete := cmdBuilder(RunImagesDelete, "delete", "Delete image", out)
	cmd.AddCommand(cmdImagesDelete)
	addStringFlag(cmdImagesDelete, doctl.ArgImageID, "", "Image id")

	return cmd
}

// RunImagesList images.
func RunImagesList(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	return listImages(ns, out, client.Images.List)
}

// RunImagesListDistribution lists distributions that are available.
func RunImagesListDistribution(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	return listImages(ns, out, client.Images.ListDistribution)
}

// RunImagesListApplication lists application iamges.
func RunImagesListApplication(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	return listImages(ns, out, client.Images.ListApplication)
}

// RunImagesListUser lists user images.
func RunImagesListUser(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	return listImages(ns, out, client.Images.ListUser)
}

// RunImagesGet retrieves an image by id or slug.
func RunImagesGet(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	rawID := doctl.DoctlConfig.GetString(ns, doctl.ArgImage)

	var err error
	var image *godo.Image
	if id, cerr := strconv.Atoi(rawID); cerr == nil {
		image, _, err = client.Images.GetByID(id)
	} else {
		if len(rawID) > 0 {
			image, _, err = client.Images.GetBySlug(rawID)
		} else {
			err = fmt.Errorf("image identifier is required")
		}
	}

	if err != nil {
		return err
	}

	return doctl.DisplayOutput(image, out)
}

// RunImagesUpdate updates an image.
func RunImagesUpdate(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetInt(ns, doctl.ArgImageID)

	req := &godo.ImageUpdateRequest{
		Name: doctl.DoctlConfig.GetString(ns, doctl.ArgImageName),
	}

	image, _, err := client.Images.Update(id, req)
	if err != nil {
		return err
	}

	return doctl.DisplayOutput(image, out)
}

// RunImagesDelete deletes an image.
func RunImagesDelete(ns string, out io.Writer) error {
	client := doctl.DoctlConfig.GetGodoClient()
	id := doctl.DoctlConfig.GetInt(ns, doctl.ArgImageID)

	_, err := client.Images.Delete(id)
	return err
}

type listFn func(*godo.ListOptions) ([]godo.Image, *godo.Response, error)

func listImages(ns string, out io.Writer, lFn listFn) error {
	public := doctl.DoctlConfig.GetBool(ns, doctl.ArgImagePublic)

	fn := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := lFn(opt)
		if err != nil {
			return nil, nil, err
		}

		si := []interface{}{}
		for _, i := range list {
			if (public && i.Public) || !public {
				si = append(si, i)
			}
		}

		return si, resp, err
	}

	si, err := doctl.PaginateResp(fn)
	if err != nil {
		return err
	}

	list := make([]godo.Image, len(si))
	for i := range si {
		list[i] = si[i].(godo.Image)
	}

	return doctl.DisplayOutput(list, out)
}

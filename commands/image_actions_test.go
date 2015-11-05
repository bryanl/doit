package commands

import (
	"io/ioutil"
	"testing"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
)

func TestImageActionsGet(t *testing.T) {
	client := &godo.Client{
		ImageActions: &doctl.ImageActionsServiceMock{
			GetFn: func(imageID, actionID int) (*godo.Action, *godo.Response, error) {
				assert.Equal(t, imageID, 1)
				assert.Equal(t, actionID, 2)
				return &testAction, nil, nil
			},
		},
	}

	withTestClient(client, func(c *TestConfig) {
		ns := "test"
		c.Set(ns, doctl.ArgImageID, 1)
		c.Set(ns, doctl.ArgActionID, 2)

		RunImageActionsGet(ns, ioutil.Discard)
	})

}

func TestImageActionsTransfer(t *testing.T) {
	client := &godo.Client{
		ImageActions: &doctl.ImageActionsServiceMock{
			TransferFn: func(imageID int, req *godo.ActionRequest) (*godo.Action, *godo.Response, error) {
				assert.Equal(t, imageID, 1)

				region := (*req)["region"]
				assert.Equal(t, region, "dev0")

				return &testAction, nil, nil
			},
		},
	}

	withTestClient(client, func(c *TestConfig) {
		ns := "test"
		c.Set(ns, doctl.ArgImageID, 1)
		c.Set(ns, doctl.ArgRegionSlug, "dev0")

		RunImageActionsTransfer(ns, ioutil.Discard)
	})
}

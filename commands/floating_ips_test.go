package commands

import (
	"io/ioutil"
	"testing"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
)

func TestFloatingIPsList(t *testing.T) {
	didRun := false

	client := &godo.Client{
		FloatingIPs: &doctl.FloatingIPsServiceMock{
			ListFn: func(opts *godo.ListOptions) ([]godo.FloatingIP, *godo.Response, error) {
				didRun = true

				resp := &godo.Response{
					Links: &godo.Links{
						Pages: &godo.Pages{},
					},
				}
				return testFloatingIPList, resp, nil
			},
		},
	}

	withTestClient(client, func(c *TestConfig) {
		ns := "test"
		RunFloatingIPList(ns, ioutil.Discard)
		if !didRun {
			t.Errorf("List() did not run")
		}
	})
}

func TestFloatingIPsGet(t *testing.T) {
	client := &godo.Client{
		FloatingIPs: &doctl.FloatingIPsServiceMock{
			GetFn: func(ip string) (*godo.FloatingIP, *godo.Response, error) {
				assert.Equal(t, "127.0.0.1", ip)
				return &testFloatingIP, nil, nil
			},
		},
	}

	withTestClient(client, func(c *TestConfig) {
		ns := "test"
		c.Set(ns, doctl.ArgIPAddress, "127.0.0.1")

		RunFloatingIPGet(ns, ioutil.Discard)
	})
}

func TestFloatingIPsCreate(t *testing.T) {
	client := &godo.Client{
		FloatingIPs: &doctl.FloatingIPsServiceMock{
			CreateFn: func(req *godo.FloatingIPCreateRequest) (*godo.FloatingIP, *godo.Response, error) {
				assert.Equal(t, "dev0", req.Region)
				assert.Equal(t, 1, req.DropletID)
				return &testFloatingIP, nil, nil
			},
		},
	}

	withTestClient(client, func(c *TestConfig) {
		ns := "test"
		c.Set(ns, doctl.ArgRegionSlug, "dev0")
		c.Set(ns, doctl.ArgDropletID, 1)

		RunFloatingIPCreate(ns, ioutil.Discard)
	})
}

func TestFloatingIPsDelete(t *testing.T) {
	client := &godo.Client{
		FloatingIPs: &doctl.FloatingIPsServiceMock{
			DeleteFn: func(ip string) (*godo.Response, error) {
				assert.Equal(t, "127.0.0.1", ip)
				return nil, nil
			},
		},
	}

	withTestClient(client, func(c *TestConfig) {
		ns := "test"
		c.Set(ns, doctl.ArgIPAddress, "127.0.0.1")

		RunFloatingIPDelete(ns, ioutil.Discard)
	})
}

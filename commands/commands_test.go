package commands

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/digitalocean/doctl"
	"github.com/digitalocean/godo"
	"github.com/spf13/viper"
)

var (
	testDroplet = godo.Droplet{
		ID: 1,
		Image: &godo.Image{
			ID:           1,
			Name:         "an-image",
			Distribution: "DOOS",
		},
		Name: "a-droplet",
		Networks: &godo.Networks{
			V4: []godo.NetworkV4{
				{IPAddress: "8.8.8.8", Type: "public"},
				{IPAddress: "172.16.1.2", Type: "private"},
			},
		},
		Region: &godo.Region{
			Slug: "test0",
			Name: "test 0",
		},
	}
	testDropletList = []godo.Droplet{testDroplet}
	testKernel      = godo.Kernel{ID: 1}
	testKernelList  = []godo.Kernel{testKernel}
	testFloatingIP  = godo.FloatingIP{
		Droplet: &testDroplet,
		Region:  testDroplet.Region,
		IP:      "127.0.0.1",
	}
	testFloatingIPList = []godo.FloatingIP{testFloatingIP}
)

type testFn func(c *TestConfig)

func withTestClient(client *godo.Client, tFn testFn) {
	ogConfig := doctl.DoctlConfig
	defer func() {
		doctl.DoctlConfig = ogConfig
	}()

	cfg := NewTestConfig(client)
	doctl.DoctlConfig = cfg

	tFn(cfg)
}

type TestConfig struct {
	Client *godo.Client
	SSHFn  func(user, host, keyPath string, port int) doctl.Runner
	v      *viper.Viper
}

func NewTestConfig(client *godo.Client) *TestConfig {
	return &TestConfig{
		Client: client,
		SSHFn: func(u, h, kp string, p int) doctl.Runner {
			logrus.WithFields(logrus.Fields{
				"user": u,
				"host": h,
			}).Info("ssh")
			return &doctl.MockRunner{}
		},
		v: viper.New(),
	}
}

var _ doctl.Config = &TestConfig{}

func (c *TestConfig) GetGodoClient() *godo.Client {
	return c.Client
}

func (c *TestConfig) SSH(user, host, keyPath string, port int) doctl.Runner {
	return c.SSHFn(user, host, keyPath, port)
}

func (c *TestConfig) Set(ns, key string, val interface{}) {
	nskey := fmt.Sprintf("%s-%s", ns, key)
	c.v.Set(nskey, val)
}

func (c *TestConfig) GetString(ns, key string) string {
	nskey := fmt.Sprintf("%s-%s", ns, key)
	return c.v.GetString(nskey)
}

func (c *TestConfig) GetInt(ns, key string) int {
	nskey := fmt.Sprintf("%s-%s", ns, key)
	return c.v.GetInt(nskey)
}

func (c *TestConfig) GetStringSlice(ns, key string) []string {
	nskey := fmt.Sprintf("%s-%s", ns, key)
	return c.v.GetStringSlice(nskey)
}

func (c *TestConfig) GetBool(ns, key string) bool {
	nskey := fmt.Sprintf("%s-%s", ns, key)
	return c.v.GetBool(nskey)
}

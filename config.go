package doit

import (
	"io"
	"io/ioutil"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// DConfig is a configuration for Doit.
type DConfig interface {
	GodoClient() *godo.Client
	Writer() io.Writer
}

type CLIConfig struct {
	token  string
	writer io.Writer
}

var _ DConfig = &CLIConfig{}

func NewCLIConfig(token string, writer io.Writer) *CLIConfig {
	return &CLIConfig{
		token:  token,
		writer: writer,
	}
}

func (c *CLIConfig) GodoClient() *godo.Client {
	ts := &TokenSource{
		AccessToken: c.token,
	}

	oc := oauth2.NewClient(oauth2.NoContext, ts)
	return godo.NewClient(oc)
}

func (c *CLIConfig) Writer() io.Writer {
	return c.writer
}

type MockConfig struct {
	GodoClientFn func() *godo.Client
	WriterFn     func() io.Writer
}

var _ DConfig = &MockConfig{}

func NewMockConfig() *MockConfig {
	return &MockConfig{
		GodoClientFn: func() *godo.Client {
			return nil
		},
		WriterFn: func() io.Writer {
			return ioutil.Discard
		},
	}
}

func (c *MockConfig) GodoClient() *godo.Client {
	return c.GodoClientFn()
}

func (c *MockConfig) Writer() io.Writer {
	return c.WriterFn()
}

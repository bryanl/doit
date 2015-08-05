package doit

import (
	"errors"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/assert"
)

var testAccount = &godo.Account{
	DropletLimit:  10,
	Email:         "user@example.com",
	UUID:          "1234",
	EmailVerified: true,
}

func TestAccountAction(t *testing.T) {
	c := NewMockConfig()
	c.GodoClientFn = func() *godo.Client {
		return &godo.Client{
			Account: &AccountServiceMock{
				GetFn: func() (*godo.Account, *godo.Response, error) {
					return testAccount, nil, nil
				},
			},
		}
	}

	a, err := AccountGet(c)
	assert.NoError(t, err)
	assert.Equal(t, testAccount, a)
}

func TestAccountAction_Error(t *testing.T) {
	c := NewMockConfig()
	c.GodoClientFn = func() *godo.Client {
		return &godo.Client{
			Account: &AccountServiceMock{
				GetFn: func() (*godo.Account, *godo.Response, error) {
					return nil, nil, errors.New("error")
				},
			},
		}
	}

	a, err := AccountGet(c)
	assert.Nil(t, a)
	assert.Error(t, err)
}

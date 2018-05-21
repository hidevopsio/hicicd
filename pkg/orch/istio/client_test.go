package istio

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)



func TestGet(t *testing.T) {
	config, err := NewClient()
	assert.Equal(t, nil, err)
	client := Client{
		Name:      "demo-provider",
		Namespace: "demo-dev",
	}
	client.Crd = config
	typ := RouterType
	c, flag := client.Get(typ)
	log.Info("get rule :", c)
	resoureceVersion := c.ResourceVersion
	assert.Equal(t, true, flag)
	log.Info("routerule get version :", resoureceVersion)
}

func TestDelete(t *testing.T) {
	config, err := NewClient()
	assert.Equal(t, nil, err)
	client := Client{
		Name:      "hello-world",
		Namespace: "demo-dev",
	}
	typ := RouterType
	client.Crd = config
	err = client.Delete(typ)
	assert.Equal(t, nil, err)
}

package istio

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)



func TestGet(t *testing.T) {
	config, err := NewClient()
	assert.Equal(t, nil, err)
	rule := Client{
		Name:      "demo-provider",
		Namespace: "demo-dev",
	}
	rule.Crd = config
	typ := RouterType
	con, flag := rule.Get(typ)
	log.Info("get rule :", con)
	resoureceVersion := con.ResourceVersion
	assert.Equal(t, true, flag)
	log.Info("route Rule get version :", resoureceVersion)
}

func TestDelete(t *testing.T) {
	config, err := NewClient()
	assert.Equal(t, nil, err)
	rule := Client{
		Name:      "demo-provider",
		Namespace: "demo-dev",
	}
	typ := RouterType
	rule.Crd = config
	err = rule.Delete(typ)
	assert.Equal(t, nil, err)
}

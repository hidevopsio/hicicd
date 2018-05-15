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
		Name:      "handler",
		Namespace: "istio-system",
	}
	rule.Crd = config
	typ :=     "service-role-binding"
	con, flag := rule.Get(typ)
	log.Info("get rule :", con)
	resoureceVersion := con.ResourceVersion
	assert.Equal(t, true, flag)
	log.Info("Roter Rule get resourece version :", resoureceVersion)
}

/*func TestUpdate(t *testing.T) {
	rule := &RouterRule{
		Client: Client{
			Labels:          map[string]string{"label": "value"},
			Annotations:     map[string]string{"annotation": "value"},
			Name:            "hello-world",
			Namespace:       "demo-dev",
			ResourceVersion: "3944181",
		},
		Route: []*routing.DestinationWeight{
			{Weight: 20, Labels: map[string]string{"version": "v1"}},
			{Weight: 80, Labels: map[string]string{"version": "v2"}},
		},
	}
	config, err := NewClient()
	rule.crd = config
	assert.Equal(t, nil, err)
	resourceVersion, err := rule.Update()
	assert.Equal(t, nil, err)
	log.Info("update rule:", resourceVersion)
}*/

func TestDelete(t *testing.T) {
	config, err := NewClient()
	assert.Equal(t, nil, err)
	rule := Client{
		Name:      "mysql-dev",
		Namespace: "istio-system",
	}
	typ := EgressType
	rule.Crd = config
	err = rule.Delete(typ)
	assert.Equal(t, nil, err)
}

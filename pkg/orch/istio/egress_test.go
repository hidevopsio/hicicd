package istio

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func TestCreateEgress(t *testing.T) {
	rule := &Egress{
		Client: Client{
			Name:        "mysql-dev",
			Namespace:   "istio-system",
		},
		Protocol:    EgressProtocol[4],
		Port:        3306,
		Destination: "172.16.8.80",
	}
	config, err := NewClient()
	rule.Crd = config
	assert.Equal(t, nil, err)
	resourceVersion, err := rule.Create()
	assert.Equal(t, nil, err)
	log.Info("create rule:", resourceVersion)
}

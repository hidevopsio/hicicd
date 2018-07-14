package istio

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func TestCreateQuota(t *testing.T) {
	rule := &Memquota{
		Client: Client{
			Labels:      map[string]string{"label": "value"},
			Annotations: map[string]string{"annotation": "value"},
			Name:        "demo-provider",
			Namespace:   "demo-dev",
		},
	}
	config, err := NewClient()
	rule.Crd = config
	assert.Equal(t, nil, err)
	resourceVersion, err := rule.Create()
	assert.Equal(t, nil, err)
	log.Info("create rule:", resourceVersion)
}


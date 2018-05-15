package istio

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	routing "istio.io/api/routing/v1alpha1"
)

func TestCreate(t *testing.T) {
	rule := &RouterRule{
		Client: Client{
			Labels:      map[string]string{"label": "value"},
			Annotations: map[string]string{"annotation": "value"},
			Name:        "demo-provider",
			Namespace:   "demo-dev",
		},

		Route: []*routing.DestinationWeight{
			{Weight: 0, Labels: map[string]string{"version": "v2"}},
			{Weight: 100, Labels: map[string]string{"version": "v1"}},
		},
		Percent:    100,
		FixedDelay: 1,
		Timeout:    1,
	}
	config, err := NewClient()
	rule.Crd = config
	assert.Equal(t, nil, err)
	resourceVersion, err := rule.Create()
	assert.Equal(t, nil, err)
	log.Info("create rule:", resourceVersion)
}

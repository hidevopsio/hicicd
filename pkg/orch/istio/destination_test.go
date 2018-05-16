package istio

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func TestCreateBreaker(t *testing.T) {
	rule := &Destination{
		Client: Client{
			Name:        "demo-provider",
			Namespace:   "demo-dev",
			Version:     "v1",
		},
		MaxConnections:               1,
		HttpMaxPendingRequests:       1,
		SleepWindow:                  1,
		HttpDetectionInterval:        1,
		HttpMaxEjectionPercent:       1,
		HttpConsecutiveErrors:        1,
		HttpMaxRequestsPerConnection: 1,
	}
	config, err := NewClient()
	rule.Crd = config
	assert.Equal(t, nil, err)
	resourceVersion, err := rule.Create()
	assert.Equal(t, nil, err)
	log.Info("create rule:", resourceVersion)
}

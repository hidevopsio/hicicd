package istio

import (
	"istio.io/istio/pilot/pkg/model"
	routing "istio.io/api/routing/v1alpha1"
	google_protobuf1 "github.com/golang/protobuf/ptypes/duration"
	"github.com/hidevopsio/hiboot/pkg/log"
)

const (
	BreakType    = "destination-policy"
	BreakVersion = "v1alpha2"
	BreakGroup   = "config.istio.io"
	BreakDomain  = "cluster"
)

type Destination struct {
	Client
	MaxConnections               int32
	HttpMaxPendingRequests       int32
	SleepWindow                  int64
	HttpDetectionInterval        int64
	HttpMaxEjectionPercent       int32
	HttpConsecutiveErrors        int32
	HttpMaxRequestsPerConnection int32
}

func (b *Destination) getConfig() (model.Config, error) {
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:        BreakType,
			Version:     BreakVersion,
			Group:       BreakGroup,
			Name:        b.Name,
			Namespace:   b.Namespace,
			Domain:      BreakDomain,
			Labels:      b.Labels,
			Annotations: b.Annotations,
		},
		Spec: &routing.DestinationPolicy{
			Destination: &routing.IstioService{
				Name:   b.Name,
				Labels: map[string]string{"version": b.Version},
			},
			CircuitBreaker: &routing.CircuitBreaker{
				CbPolicy: &routing.CircuitBreaker_SimpleCb{
					SimpleCb: &routing.CircuitBreaker_SimpleCircuitBreakerPolicy{
						MaxConnections:         b.MaxConnections,
						HttpMaxPendingRequests: b.HttpMaxPendingRequests,
						SleepWindow: &google_protobuf1.Duration{
							Seconds: b.SleepWindow,
						},
						HttpDetectionInterval: &google_protobuf1.Duration{
							Seconds: b.HttpDetectionInterval,
						},
						HttpMaxEjectionPercent:       b.HttpMaxEjectionPercent,
						HttpConsecutiveErrors:        b.HttpConsecutiveErrors,
						HttpMaxRequestsPerConnection: b.HttpMaxRequestsPerConnection,
					},
				},
			},
		},
	}
	return config, nil
}

func (b *Destination) Create() (string, error) {
	log.Debug("create rule :", b)
	config, err := b.getConfig()
	con, exists := b.Get(BreakType)
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := b.Crd.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create route rule config ", config)
	resourceVersion, err := b.Crd.Create(config)
	if err != nil {
		log.Error("create route rule error %v", err)
		return "", err
	}
	return resourceVersion, nil
}

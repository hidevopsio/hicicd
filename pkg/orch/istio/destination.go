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
	MaxConnections               int32 `json:"max_connections"`
	HttpMaxPendingRequests       int32 `json:"http_max_pending_requests"`
	SleepWindow                  int64 `json:"sleep_window"`
	HttpDetectionInterval        int64 `json:"http_detection_interval"`
	HttpMaxEjectionPercent       int32 `json:"http_max_ejection_percent"`
	HttpConsecutiveErrors        int32 `json:"http_consecutive_errors"`
	HttpMaxRequestsPerConnection int32 `json:"http_max_requests_per_connection"`
}

func (d *Destination) getConfig() (model.Config, error) {
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:        BreakType,
			Version:     BreakVersion,
			Group:       BreakGroup,
			Name:        d.Name,
			Namespace:   d.Namespace,
			Domain:      BreakDomain,
			Labels:      d.Labels,
			Annotations: d.Annotations,
		},
		Spec: &routing.DestinationPolicy{
			Destination: &routing.IstioService{
				Name:   d.Name,
				Labels: map[string]string{"version": d.Version},
			},
			CircuitBreaker: &routing.CircuitBreaker{
				CbPolicy: &routing.CircuitBreaker_SimpleCb{
					SimpleCb: &routing.CircuitBreaker_SimpleCircuitBreakerPolicy{
						MaxConnections:         d.MaxConnections,
						HttpMaxPendingRequests: d.HttpMaxPendingRequests,
						SleepWindow: &google_protobuf1.Duration{
							Seconds: d.SleepWindow,
						},
						HttpDetectionInterval: &google_protobuf1.Duration{
							Seconds: d.HttpDetectionInterval,
						},
						HttpMaxEjectionPercent:       d.HttpMaxEjectionPercent,
						HttpConsecutiveErrors:        d.HttpConsecutiveErrors,
						HttpMaxRequestsPerConnection: d.HttpMaxRequestsPerConnection,
					},
				},
			},
		},
	}
	return config, nil
}

func (d *Destination) Create() (string, error) {
	log.Debug("create destination :", d)
	config, err := d.getConfig()
	con, exists := d.Get(BreakType)
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := d.Crd.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create destination ", config)
	resourceVersion, err := d.Crd.Create(config)
	if err != nil {
		log.Error("create destination error :", err)
		return "", err
	}
	return resourceVersion, nil
}

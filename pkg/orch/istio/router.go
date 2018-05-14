package istio

import (
	"istio.io/istio/pilot/pkg/model"
	routing "istio.io/api/routing/v1alpha1"
	"github.com/hidevopsio/hiboot/pkg/log"
	google_protobuf1 "github.com/golang/protobuf/ptypes/duration"
)

const (
	RouterType    = "route-rule"
	RouterVersion = "v1alpha2"
	RouterGroup   = "config.istio.io"
	RouterDomain  = "cluster"
)

type RouterRule struct {
	Client
	Route          []*routing.DestinationWeight
	HttpFault      *routing.HTTPFaultInjection
	HttpReqTimeout *routing.HTTPTimeout
	Timeout        int64
	Percent        float32
	FixedDelay     int64
}

func (router *RouterRule) getConfig() (model.Config, error) {
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:        RouterType,
			Version:     RouterVersion,
			Group:       RouterGroup,
			Name:        router.Name,
			Namespace:   router.Namespace,
			Domain:      RouterDomain,
			Labels:      router.Labels,
			Annotations: router.Annotations,
		},
		Spec: &routing.RouteRule{
			Destination: &routing.IstioService{
				Name: router.Name,
			},
			HttpFault: &routing.HTTPFaultInjection{
				Delay: &routing.HTTPFaultInjection_Delay{
					Percent: router.Percent,
					HttpDelayType: &routing.HTTPFaultInjection_Delay_FixedDelay{
						FixedDelay: &google_protobuf1.Duration{
							Seconds: router.FixedDelay,
						},
					},
				},
			},
			HttpReqTimeout: &routing.HTTPTimeout{
				TimeoutPolicy: &routing.HTTPTimeout_SimpleTimeout{
					SimpleTimeout: &routing.HTTPTimeout_SimpleTimeoutPolicy{
						Timeout: &google_protobuf1.Duration{
							Seconds: router.Timeout,
						},
					},
				},
			},
			Route: router.Route,
		},
	}
	return config, nil
}

func (router *RouterRule) Create() (string, error) {
	log.Debug("create rule :", router)
	config, err := router.getConfig()
	con, exists := router.Get(RouterType)
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := router.crd.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create route rule config ", config)
	resourceVersion, err := router.crd.Create(config)
	if err != nil {
		log.Error("create route rule error %v", err)
		return "", err
	}
	return resourceVersion, nil
}

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
	Route          []*routing.DestinationWeight `json:"route"`
	HttpFault      *routing.HTTPFaultInjection  `json:"http_fault"`
	HttpReqTimeout *routing.HTTPTimeout         `json:"http_req_timeout"`
	Timeout        int64                        `json:"timeout"`
	Percent        float32                      `json:"percent"`
	FixedDelay     int64                        `json:"fixed_delay"`
}

func (router *RouterRule) GetConfig() (model.Config, error) {
	route := &routing.RouteRule{
		Destination: &routing.IstioService{
			Name: router.Name,
		},
	}
	if router.FixedDelay != 0 {
		route.HttpFault = &routing.HTTPFaultInjection{
			Delay: &routing.HTTPFaultInjection_Delay{
				Percent: router.Percent,
				HttpDelayType: &routing.HTTPFaultInjection_Delay_FixedDelay{
					FixedDelay: &google_protobuf1.Duration{
						Seconds: router.FixedDelay,
					},
				},
			},
		}
	}
	if router.Timeout != 0 {
		route.HttpReqTimeout = &routing.HTTPTimeout{
			TimeoutPolicy: &routing.HTTPTimeout_SimpleTimeout{
				SimpleTimeout: &routing.HTTPTimeout_SimpleTimeoutPolicy{
					Timeout: &google_protobuf1.Duration{
						Seconds: router.Timeout,
					},
				},
			},
		}
	}
	if len(router.Route) != 0 {
		route.Route = router.Route
	}
	log.Debug("route.HttpFault :{}", route)
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:        RouterType,
			Version:     RouterVersion,
			Group:       RouterGroup,
			Name:        router.Name,
			Namespace:   router.Namespace,
			Domain:      RouterDomain,
			Labels:      map[string]string{"label": "value"},
			Annotations: map[string]string{"annotation": "value"},
		},
		Spec: route,
	}
	return config, nil
}

func (router *RouterRule) Create() (string, error) {
	log.Debug("create rule :", router)
	config, err := router.GetConfig()
	con, exists := router.Get(RouterType)
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := router.Crd.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create route rule config ", config)
	resourceVersion, err := router.Crd.Create(config)
	if err != nil {
		log.Error("create route rule error %v", err)
		return "", err
	}
	return resourceVersion, nil
}

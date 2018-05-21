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
//Obtain  config
func (rr *RouterRule) GetConfig() (model.Config, error) {
	routerule := &routing.RouteRule{
		Destination: &routing.IstioService{
			Name: rr.Name,
		},
	}
	if rr.FixedDelay != 0 {
		routerule.HttpFault = &routing.HTTPFaultInjection{
			Delay: &routing.HTTPFaultInjection_Delay{
				Percent: rr.Percent,
				HttpDelayType: &routing.HTTPFaultInjection_Delay_FixedDelay{
					FixedDelay: &google_protobuf1.Duration{
						Seconds: rr.FixedDelay,
					},
				},
			},
		}
	}
	if rr.Timeout != 0 {
		routerule.HttpReqTimeout = &routing.HTTPTimeout{
			TimeoutPolicy: &routing.HTTPTimeout_SimpleTimeout{
				SimpleTimeout: &routing.HTTPTimeout_SimpleTimeoutPolicy{
					Timeout: &google_protobuf1.Duration{
						Seconds: rr.Timeout,
					},
				},
			},
		}
	}
	if len(rr.Route) != 0 {
		routerule.Route = rr.Route
	}
	log.Debug("route.HttpFault :{}", routerule)
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:        RouterType,
			Version:     RouterVersion,
			Group:       RouterGroup,
			Name:        rr.Name,
			Namespace:   rr.Namespace,
			Domain:      RouterDomain,
			Labels:      map[string]string{"label": "value"},
			Annotations: map[string]string{"annotation": "value"},
		},
		Spec: routerule,
	}
	return config, nil
}
//create  route rule
func (rr *RouterRule) Create() (string, error) {
	log.Debug("create rule :", rr)
	config, err := rr.GetConfig()
	con, exists := rr.Get(RouterType)
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := rr.Crd.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create route rule config ", config)
	resourceVersion, err := rr.Crd.Create(config)
	if err != nil {
		log.Error("create route rule error %v", err)
		return "", err
	}
	return resourceVersion, nil
}

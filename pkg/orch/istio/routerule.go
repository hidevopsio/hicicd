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
func (r *RouterRule) GetConfig() (model.Config, error) {
	routerule := &routing.RouteRule{
		Destination: &routing.IstioService{
			Name: r.Name,
		},
	}
	if r.FixedDelay != 0 {
		routerule.HttpFault = &routing.HTTPFaultInjection{
			Delay: &routing.HTTPFaultInjection_Delay{
				Percent: r.Percent,
				HttpDelayType: &routing.HTTPFaultInjection_Delay_FixedDelay{
					FixedDelay: &google_protobuf1.Duration{
						Seconds: r.FixedDelay,
					},
				},
			},
		}
	}
	if r.Timeout != 0 {
		routerule.HttpReqTimeout = &routing.HTTPTimeout{
			TimeoutPolicy: &routing.HTTPTimeout_SimpleTimeout{
				SimpleTimeout: &routing.HTTPTimeout_SimpleTimeoutPolicy{
					Timeout: &google_protobuf1.Duration{
						Seconds: r.Timeout,
					},
				},
			},
		}
	}
	if len(r.Route) != 0 {
		routerule.Route = r.Route
	}
	log.Debug("route.HttpFault :{}", routerule)
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:        RouterType,
			Version:     RouterVersion,
			Group:       RouterGroup,
			Name:        r.Name,
			Namespace:   r.Namespace,
			Domain:      RouterDomain,
			Labels:      map[string]string{"label": "value"},
			Annotations: map[string]string{"annotation": "value"},
		},
		Spec: routerule,
	}
	return config, nil
}
//create  route rule
func (r *RouterRule) Create() (string, error) {
	log.Debug("create rule :", r)
	config, err := r.GetConfig()
	con, exists := r.Get(RouterType)
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := r.Crd.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create route rule config ", config)
	resourceVersion, err := r.Crd.Create(config)
	if err != nil {
		log.Error("create route rule error %v", err)
		return "", err
	}
	return resourceVersion, nil
}

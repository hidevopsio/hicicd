package istio

import "istio.io/istio/pilot/pkg/model"
import (
	routing "istio.io/api/routing/v1alpha1"
	"github.com/hidevopsio/hiboot/pkg/log"
)

const (
	EgressType    = "egress-rule"
	EgressVersion = "v1alpha2"
	EgressGroup   = "config.istio.io"
	EgressDomain  = "cluster"
)

//HTTP|HTTPS|GRPC|HTTP2|TCP|MONGO
var EgressProtocol = [...]string{"HTTP", "HTTPS", "GRPC", "HTTP2", "TCP", "MONGO"}

type Egress struct {
	Client
	Destination string
	Port        int32
	Protocol    string
}

func (e *Egress) getConfig() (model.Config, error) {
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:        EgressType,
			Version:     EgressVersion,
			Group:       EgressGroup,
			Name:        e.Name,
			Namespace:   e.Namespace,
			Domain:      EgressDomain,
			Labels:      e.Labels,
			Annotations: e.Annotations,
		},
		Spec: &routing.EgressRule{
			Destination: &routing.IstioService{
				Service: e.Destination,
			},
			Ports: []*routing.EgressRule_Port{
				{
					Port:     e.Port,
					Protocol: e.Protocol,
				},
			},
		},
	}
	return config, nil
}

func (egress *Egress) Create() (string, error) {
	log.Debug("create rule :", egress)
	config, err := egress.getConfig()
	con, exists := egress.Get(EgressType)
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := egress.crd.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create route rule config ", config)
	resourceVersion, err := egress.crd.Create(config)
	if err != nil {
		log.Error("create route rule error %v", err)
		return "", err
	}
	return resourceVersion, nil
}

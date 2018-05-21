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
	Destination string `json:"destination"`
	Port        int32  `json:"port"`
	Protocol    string `json:"protocol"`
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

func (e *Egress) Create() (string, error) {
	log.Debug("create egress :", e)
	config, err := e.getConfig()
	con, exists := e.Get(EgressType)
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := e.Crd.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create egress config ", config)
	resourceVersion, err := e.Crd.Create(config)
	if err != nil {
		log.Error("create egress error %v", err)
		return "", err
	}
	return resourceVersion, nil
}

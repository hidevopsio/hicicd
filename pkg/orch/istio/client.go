package istio

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"github.com/hidevopsio/hicicd/pkg/orch"
)

type ClientInterface interface {
	Create() (string, error)
	Get(typ string) (*model.Config, bool)
	Delete(typ string) error
	Update() (string, error)
	getConfig() (model.Config, error)
}

type Client struct {
	Type            int64 `json:"type"`
	FullName        string `json:"full_name"`
	Version         string `json:"version"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	Group           string `json:"group"`
	Labels          map[string]string `json:"labels"`
	Annotations     map[string]string `json:"annotations"`
	Domain          string `json:"domain"`
	ResourceVersion string `json:"resource_version"`
	Crd             *crd.Client
}

var Typ = [5]string{"route-rule","egress-rule","destination-policy","quota-spec"}

const (
	Type    = "route-rule"
	Version = "v1alpha2"
	Group   = "config.istio.io"
	Domain  = "cluster"
)

func newClient(kubeconfig string) (*crd.Client, error) {
	// TODO: use model.IstioConfigTypes once model.IngressRule is deprecated
	config := model.ConfigDescriptor{
		model.RouteRule,
		model.VirtualService,
		model.Gateway,
		model.EgressRule,
		model.ServiceEntry,
		model.DestinationPolicy,
		model.DestinationRule,
		model.HTTPAPISpec,
		model.HTTPAPISpecBinding,
		model.QuotaSpec,
		model.QuotaSpecBinding,
		model.AuthenticationPolicy,
		model.ServiceRole,
		model.ServiceRoleBinding,
	}
	return crd.NewClient(kubeconfig, config, "")
}

func (client *Client) getConfig() (*model.Config, error) {
	return nil, nil
}

func (c *Client) Get(typ string) (*model.Config, bool){
	config, exists := c.Crd.Get(typ, c.Name, c.Namespace)
	log.Debug("route rule get config", exists)
	return config, exists
}

func (client *Client) Delete(typ string) error {
	err := client.Crd.Delete(typ, client.Name, client.Namespace)
	if err != nil {
		log.Error("type: "+ typ +" route rule delete config", err)
		return err
	}
	return nil
}

func NewClient() (*crd.Client, error) {
	log.Debug("create config kubeconfig")
	configClient, err := newClient(*orch.Kubeconfig)
	if err != nil {
		log.Error("create config configClient error", err)
		return nil, err
	}
	return configClient, nil
}

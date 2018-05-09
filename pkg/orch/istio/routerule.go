package istio

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	routing "istio.io/api/routing/v1alpha1"
	"github.com/hidevopsio/hicicd/pkg/orch"
)

type Client struct {
	Type            string
	Version         string
	Name            string
	Namespace       string
	Group           string
	Labels          map[string]string
	Annotations     map[string]string
	Domain          string
	ResourceVersion string
	Route           []*routing.DestinationWeight
	crd             *crd.Client
}

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

func (client *Client) Create() (string, error) {
	log.Debug("create rule :", client)
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:        Type,
			Version:     Version,
			Group:       Group,
			Name:        client.Name,
			Namespace:   client.Namespace,
			Domain:      Domain,
			Labels:      client.Labels,
			Annotations: client.Annotations,
		},
		Spec: &routing.RouteRule{
			Destination: &routing.IstioService{
				Name: "reviews",
			},
			Route: client.Route,
		},
	}
	con, exists := client.Get()
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := client.crd.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create route rule config ", config)
	resourceVersion, err := client.crd.Create(config)
	if err != nil {
		log.Error("create route rule error %v", err)
		return "", err
	}
	return resourceVersion, nil
}

func (client *Client) Get() (*model.Config, bool) {
	config, flag := client.crd.Get(Type, client.Name, client.Namespace)
	log.Debug("route rule get config", flag)
	return config, flag
}

func (client *Client) Delete() error {
	err := client.crd.Delete(Type, client.Name, client.Namespace)
	if err != nil {
		log.Error("route rule delete config", err)
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

func (client *Client) Update() (string, error) {
	log.Debug("update rule :", client)
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:            Type,
			Version:         Version,
			Group:           Group,
			Name:            client.Name,
			Namespace:       client.Namespace,
			Domain:          Domain,
			ResourceVersion: client.ResourceVersion,
			Labels:          client.Labels,
			Annotations:     client.Annotations,
		},
		Spec: &routing.RouteRule{
			Destination: &routing.IstioService{
				Name: "reviews",
			},
			Route: client.Route,
		},
	}
	log.Debug("update route rule config ", config)
	resourceVersion, err := client.crd.Update(config)
	if err != nil {
		log.Error("update route rule error ", err)
		return "", err
	}
	return resourceVersion, nil
}

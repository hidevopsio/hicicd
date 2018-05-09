package istio

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	routing "istio.io/api/routing/v1alpha1"
)

type Rule struct {
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

func (rule *Rule) Create(configClient *crd.Client) (string, error) {
	log.Debug("create rule :", rule)
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:        Type,
			Version:     Version,
			Group:       Group,
			Name:        rule.Name,
			Namespace:   rule.Namespace,
			Domain:      Domain,
			Labels:      rule.Labels,
			Annotations: rule.Annotations,
		},
		Spec: &routing.RouteRule{
			Destination: &routing.IstioService{
				Name: "reviews",
			},
			Route: rule.Route,
		},
	}
	con, exists := rule.Get(configClient)
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := configClient.Update(config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create router rule config ", config)
	resourceVersion, err := configClient.Create(config)
	if err != nil {
		log.Error("create router rule error %v", err)
		return "", err
	}
	return resourceVersion, nil
}

func (rule *Rule) Get(configClient *crd.Client) (*model.Config, bool) {
	config, flag := configClient.Get(Type, rule.Name, rule.Namespace)
	log.Debug("roter rule get config", flag)
	return config, flag
}

func (rule *Rule) Delete(configClient *crd.Client) error {
	err := configClient.Delete(Type, rule.Name, rule.Namespace)
	if err != nil {
		log.Error("roter rule delete config", err)
		return err
	}
	return nil
}

func NewClient(kubeconfig string) (*crd.Client, error) {
	log.Debug("create config kubeconfig", kubeconfig)
	configClient, err := newClient(kubeconfig)
	if err != nil {
		log.Error("create config configClient error", err)
		return nil, err
	}
	return configClient, nil
}

func (rule *Rule) Update(configClient *crd.Client) (string, error) {
	log.Debug("update rule :", rule)
	config := model.Config{
		ConfigMeta: model.ConfigMeta{
			Type:            Type,
			Version:         Version,
			Group:           Group,
			Name:            rule.Name,
			Namespace:       rule.Namespace,
			Domain:          Domain,
			ResourceVersion: rule.ResourceVersion,
			Labels:          rule.Labels,
			Annotations:     rule.Annotations,
		},
		Spec: &routing.RouteRule{
			Destination: &routing.IstioService{
				Name: "reviews",
			},
			Route: rule.Route,
		},
	}
	log.Debug("update router rule config ", config)
	resourceVersion, err := configClient.Update(config)
	if err != nil {
		log.Error("update router rule error ", err)
		return "", err
	}
	return resourceVersion, nil
}

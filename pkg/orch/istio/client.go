package istio

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"github.com/hidevopsio/hicicd/pkg/orch"
)

type ClientInterface interface {
	Create() (string, error)
	Get() (*model.Config, bool)
	Delete() error
	Update() (string, error)
	getConfig() (model.Config, error)
}

type Client struct {
	Type            string
	FullName        string
	Version         string
	Name            string
	Namespace       string
	Group           string
	Labels          map[string]string
	Annotations     map[string]string
	Domain          string
	ResourceVersion string
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

func (client *Client) getConfig() (*model.Config, error) {
	return nil, nil
}

/*func (client *Client) Create() (string, error) {
	log.Debug("create rule :", client)
	config, err := client.getConfig()
	if config == nil {
		log.Error("client :{}", err)
		return "", err
	}
	con, exists := client.Get()
	log.Debug("config exists", exists)
	if exists {
		config.ResourceVersion = con.ResourceVersion
		resourceVersion, err := client.crd.Update(*config)
		if err != nil {
			return "", err
		}
		return resourceVersion, nil
	}
	log.Debug("create route rule config ", config)
	resourceVersion, err := client.crd.Create(*config)
	if err != nil {
		log.Error("create route rule error %v", err)
		return "", err
	}
	return resourceVersion, nil
}*/

func (client *Client) Get() (*model.Config, bool) {
	config, flag := client.crd.Get(Type, client.Name, client.Namespace)
	log.Debug("route rule get config", flag)
	return config, flag
}

func (client *Client) Delete(typ string) error {
	err := client.crd.Delete(typ, client.Name, client.Namespace)
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

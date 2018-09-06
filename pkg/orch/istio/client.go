package istio

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"github.com/hidevopsio/hicicd/pkg/orch"
	"os"
	"github.com/hashicorp/go-multierror"
	"fmt"
	"os/user"
)

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
	log.Debug("kubeconfig :", kubeconfig)
	return crd.NewClient(kubeconfig, config, "")
}

func (c *Client) getConfig() (*model.Config, error) {
	return nil, nil
}

func (c *Client) Get(typ string) (*model.Config, bool){
	config, exists := c.Crd.Get(typ, c.Name, c.Namespace)
	log.Debug("route rule get config", exists)
	return config, exists
}

func (c *Client) Delete(typ string) error {
	err := c.Crd.Delete(typ, c.Name, c.Namespace)
	if err != nil {
		log.Error("type: "+ typ +" route rule delete config", err)
		return err
	}
	return nil
}

func NewClient() (*crd.Client, error) {
	log.Debug("create config kubeconfig", *orch.Kubeconfig)
	configClient, err := newClient(*orch.Kubeconfig)
	if err != nil {
		log.Error("create config configClient error", err)
		return nil, err
	}
	return configClient, nil
}

func ResolveConfig(kubeconfig string) (string, error) {
	// Consistency with kubectl
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}
	if kubeconfig == "" {
		usr, err := user.Current()
		if err == nil {
			defaultCfg := usr.HomeDir + "/.kube/config"
			_, err := os.Stat(kubeconfig)
			if err != nil {
				kubeconfig = defaultCfg
			}
		}
	}
	if kubeconfig != "" {
		info, err := os.Stat(kubeconfig)
		if err != nil {
			if os.IsNotExist(err) {
				err = fmt.Errorf("kubernetes configuration file %q does not exist", kubeconfig)
			} else {
				err = multierror.Append(err, fmt.Errorf("kubernetes configuration file %q", kubeconfig))
			}
			return "", err
		}

		// if it's an empty file, switch to in-cluster config
		if info.Size() == 0 {
			log.Info("using in-cluster configuration")
			return "", nil
		}
	}
	return kubeconfig, nil
}

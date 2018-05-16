package openshift


import (
	policyv1 "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	"github.com/hidevopsio/hicicd/pkg/orch"
	"github.com/hidevopsio/hiboot/pkg/log"
	v1 "github.com/openshift/api/authorization/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

)


type Policy struct{
	Name string
	Namespace string

	Interface policyv1.RoleInterface
}

func NewPolicy(name, namespace string) (*Policy, error)  {
	log.Debug("NewPolicy()")
	client, err := policyv1.NewForConfig(orch.Config)
	if err != nil {
		return nil, err
	}
	r := &Policy{
		Name:      name,
		Namespace: namespace,
		Interface: client.Roles(namespace),
	}
	return r, nil
}

func (p *Policy) Create(policy *v1.Role) (*v1.Role, error){
	log.Debug("create policy:", policy)
	po, err := p.Interface.Create(policy)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (p *Policy) Get() (*v1.Role, error){
	log.Debug("get policy:")
	po, err:= p.Interface.Get(p.Name, meta_v1.GetOptions{})
	if err != nil {
		log.Error("get policy err :", err)
		return nil, err
	}
	return po, nil
}

package openshift

import (
	"github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/hidevopsio/hicicd/pkg/orch"
	"github.com/hidevopsio/hiboot/pkg/log"
	authorization_v1 "github.com/openshift/api/authorization/v1"
)

type RoleBinding struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Data      map[string]string `json:"data"`
	Interface v1.RoleBindingInterface
}

func NewRoleBinding(name, namespace string) (*RoleBinding, error) {
	log.Debug("NewPolicy()")
	client, err := v1.NewForConfig(orch.Config)
	if err != nil {
		return nil, err
	}
	r := &RoleBinding{
		Name:      name,
		Namespace: namespace,
		Interface: client.RoleBindings(namespace),
	}
	return r, nil
}

func (r *RoleBinding) Get() (*authorization_v1.RoleBinding, error) {
	log.Debug("get RoleBinding:")
	role, err := r.Interface.Get(r.Name, meta_v1.GetOptions{})
	if err != nil {
		log.Error("get policy err :", err)
		return nil, err
	}
	return role, nil
}

func (r *RoleBinding) Create(rolebinding *authorization_v1.RoleBinding) (*authorization_v1.RoleBinding, error) {
	log.Debug("create rolebinding")
	_, err := r.Interface.Get(r.Name, meta_v1.GetOptions{})
	if err != nil {
		log.Error("get policy err :", err)
		result, err := r.Update(rolebinding)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	result, err := r.Interface.Create(rolebinding)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RoleBinding) Delete() error {
	log.Debug("get RoleBinding:")
	err := r.Interface.Delete(r.Name, &meta_v1.DeleteOptions{})
	return err
}

func (r *RoleBinding) Update(rolebinding *authorization_v1.RoleBinding) (*authorization_v1.RoleBinding, error) {
	log.Debug("get RoleBinding:")
	result, err := r.Interface.Update(rolebinding)
	if err != nil {
		return nil, err
	}
	return result, nil
}

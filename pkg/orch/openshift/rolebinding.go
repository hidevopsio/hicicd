package openshift

import (
	"github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/hidevopsio/hicicd/pkg/orch"
	"github.com/hidevopsio/hiboot/pkg/log"
	authorization_v1 "github.com/openshift/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
)

type RoleBinding struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	RoleRefName string            `json:"role_ref_name"`
	RoleRefKind string            `json:"role_ref_kind"`
	SubjectKind string            `json:"subject_kind"`
	SubjectName string            `json:"subject_name"`
	Data        map[string]string `json:"data"`
	Interface   v1.RoleBindingInterface
}

func (r *RoleBinding) Init() (*authorization_v1.RoleBinding) {
	roleBinding := &authorization_v1.RoleBinding{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      r.Name,
			Namespace: r.Namespace,
		},
		RoleRef: corev1.ObjectReference{
			Name: r.RoleRefName,
			Kind: r.RoleRefKind,
		},
		Subjects: []corev1.ObjectReference{
			{
				Kind:      r.SubjectKind,
				Name:      r.SubjectName,
				Namespace: r.Namespace,
			},
		},
	}
	return roleBinding
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
	log.Debug("create role binding")
	_, err := r.Interface.Get(r.Name, meta_v1.GetOptions{})
	if err == nil {
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

func (r *RoleBinding) Update(roleBinding *authorization_v1.RoleBinding) (*authorization_v1.RoleBinding, error) {
	log.Debug("get RoleBinding:")
	result, err := r.Interface.Update(roleBinding)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RoleBinding) InitImagePullers() error {
	name := "system:image-pullers"
	namespace := r.Namespace
	bin, err := NewRoleBinding(name, namespace)
	if err != nil {
		return err
	}
	roleBinding := &authorization_v1.RoleBinding{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      name,
			Namespace: r.Namespace,
		},
		RoleRef: corev1.ObjectReference{
			Name: "system:image-puller",
			Kind: "ClusterRole",
		},
		Subjects: []corev1.ObjectReference{
			{
				Kind:      "Group",
				Name:      "system:serviceaccounts:" + namespace,
				Namespace: namespace,
			},
		},
	}
	_, err = bin.Create(roleBinding)
	return err
}

func (r *RoleBinding) InitImageBuilders() error {
	name := "system:image-builders"
	namespace := r.Namespace
	bin, err := NewRoleBinding(name, namespace)
	if err != nil {
		return err
	}
	roleBinding := &authorization_v1.RoleBinding{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		RoleRef: corev1.ObjectReference{
			Name: "system:image-builder",
			Kind: "ClusterRole",
		},
		Subjects: []corev1.ObjectReference{
			{
				Kind:      "ServiceAccount",
				Name:      "builder",
				Namespace: namespace,
			},
		},
	}
	_, err = bin.Create(roleBinding)
	return err
}

func (r *RoleBinding) InitSystemDeployers() error {
	name := "system:deployers"
	namespace := r.Namespace
	bin, err := NewRoleBinding(name, namespace)
	if err != nil {
		return err
	}
	roleBinding := &authorization_v1.RoleBinding{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		RoleRef: corev1.ObjectReference{
			Name: "system:deployer",
			Kind: "ClusterRole",
		},
		Subjects: []corev1.ObjectReference{
			{
				Kind:      "ServiceAccount",
				Name:      "deployer",
				Namespace: namespace,
			},
		},
	}
	_, err = bin.Create(roleBinding)
	return err
}
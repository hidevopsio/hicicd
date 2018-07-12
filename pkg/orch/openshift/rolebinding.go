package openshift

import (
	"github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	"github.com/openshift/client-go/authorization/clientset/versioned/fake"
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

func (rb *RoleBinding) Init() (*authorization_v1.RoleBinding) {
	roleBinding := &authorization_v1.RoleBinding{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      rb.Name,
			Namespace: rb.Namespace,
		},
		RoleRef: corev1.ObjectReference{
			Name: rb.RoleRefName,
			Kind: rb.RoleRefKind,
		},
		Subjects: []corev1.ObjectReference{
			{
				Kind:      rb.SubjectKind,
				Name:      rb.SubjectName,
				Namespace: rb.Namespace,
			},
		},
	}
	return roleBinding
}



func NewRoleBindingClientSet() (v1.AuthorizationV1Interface, error) {

	cli := orch.GetClientInstance()

	// get the fake ClientSet for testing
	if cli.IsTestRunning() {
		return fake.NewSimpleClientset().AuthorizationV1(), nil
	}

	// get the real ClientSet
	clientSet, err := v1.NewForConfig(cli.Config())

	return clientSet, err
}


func NewRoleBinding(name, namespace string) (*RoleBinding, error) {
	log.Debug("NewPolicy()")
	client, err := NewRoleBindingClientSet()
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

func (rb *RoleBinding) Get() (*authorization_v1.RoleBinding, error) {
	log.Debug("get RoleBinding:")
	role, err := rb.Interface.Get(rb.Name, meta_v1.GetOptions{})
	if err != nil {
		log.Error("get policy err :", err)
		return nil, err
	}
	return role, nil
}

func (rb *RoleBinding) Create(roleBinding *authorization_v1.RoleBinding) (*authorization_v1.RoleBinding, error) {
	log.Debug("create role binding")
	_, err := rb.Interface.Get(rb.Name, meta_v1.GetOptions{})
	if err == nil {
		result, err := rb.Update(roleBinding)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	result, err := rb.Interface.Create(roleBinding)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rb *RoleBinding) Delete() error {
	log.Debug("get RoleBinding:")
	err := rb.Interface.Delete(rb.Name, &meta_v1.DeleteOptions{})
	return err
}

func (rb *RoleBinding) Update(roleBinding *authorization_v1.RoleBinding) (*authorization_v1.RoleBinding, error) {
	log.Debug("get RoleBinding:")
	result, err := rb.Interface.Update(roleBinding)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rb *RoleBinding) InitImagePullers() error {
	name := "system:image-pullers"
	namespace := rb.Namespace
	bin, err := NewRoleBinding(name, namespace)
	if err != nil {
		return err
	}
	roleBinding := &authorization_v1.RoleBinding{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      name,
			Namespace: rb.Namespace,
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

func (rb *RoleBinding) InitImageBuilders() error {
	name := "system:image-builders"
	namespace := rb.Namespace
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

func (rb *RoleBinding) InitSystemDeployers() error {
	name := "system:deployers"
	namespace := rb.Namespace
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
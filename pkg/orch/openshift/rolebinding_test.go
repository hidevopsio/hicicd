package openshift

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	authorization_v1 "github.com/openshift/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestRoleBindingGet(t *testing.T) {
	name := "admin"
	namespace := "demo-dev"
	bin, err := NewRoleBinding(name, namespace)
	assert.Equal(t, nil, err)
	role, err := bin.Get()
	assert.Equal(t, nil, err)
	log.Debug(role)
}

func TestRoleBindingDelete(t *testing.T) {
	name := "admin"
	namespace := "default"
	bin, err := NewRoleBinding(name, namespace)
	assert.Equal(t, nil, err)
	err = bin.Delete()
	assert.Equal(t, nil, err)
}

func TestRoleBindingCreate(t *testing.T) {
	name := "admin"
	namespace := "demo-test"
	bin, err := NewRoleBinding(name, namespace)
	assert.Equal(t, nil, err)
	rb := &RoleBinding{
		Name:name,
		Namespace: namespace,
		RoleRefName: "admin",
		RoleRefKind: "",
		SubjectName: "chen",
		SubjectKind: "User",
	}
	rolebinding := rb.Init()
	role, err := bin.Create(rolebinding)
	log.Debug(role)
	log.Debug(err)
}


func TestCreateImagePullers(t *testing.T)  {
	name := "system:image-pullers"
	namespace := "demo-test"
	bin, err := NewRoleBinding(name, namespace)
	assert.Equal(t, nil, err)
	rolebinding := &authorization_v1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
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
	role, err := bin.Create(rolebinding)
	log.Debug(role)
	log.Debug(err)

}


func TestCreateImageBuilders(t *testing.T)  {

	name := "system:image-builders"
	namespace := "demo-test"
	bin, err := NewRoleBinding(name, namespace)
	assert.Equal(t, nil, err)
	rolebinding := &authorization_v1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
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
	role, err := bin.Create(rolebinding)
	log.Debug(role)
	log.Debug(err)


}



func TestCreateSystemDeployers(t *testing.T) {
	name := "system:deployers"
	namespace := "demo-test"
	bin, err := NewRoleBinding(name, namespace)
	assert.Equal(t, nil, err)
	rolebinding := &authorization_v1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
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
	role, err := bin.Create(rolebinding)
	log.Debug(role)
	log.Debug(err)

}

func TestRoleBindingUpdate(t *testing.T) {
	name := "admin"
	namespace := "default"
	bin, err := NewRoleBinding(name, namespace)
	assert.Equal(t, nil, err)
	rolebinding := &authorization_v1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		RoleRef: corev1.ObjectReference{
			Name: name,
		},
		Subjects: []corev1.ObjectReference{
			{
				Name: "chen",
				Kind: "User",
			}, {
				Kind: "User",
				Name: "shi",
			},
		},
	}
	role, err := bin.Update(rolebinding)
	log.Debug(role)
	log.Debug(err)
}

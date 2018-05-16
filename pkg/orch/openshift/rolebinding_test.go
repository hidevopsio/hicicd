package openshift

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	authorization_v1 "github.com/openshift/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestRoleBinding_Get(t *testing.T) {
	name := "admin"
	namespace := "demo-dev"
	bin, err := NewRoleBinding(name, namespace)
	assert.Equal(t, nil, err)
	role, err := bin.Get()
	assert.Equal(t, nil, err)
	log.Debug(role)
}

func TestRoleBinding_Delete(t *testing.T) {
	name := "admin"
	namespace := "default"
	bin, err := NewRoleBinding(name, namespace)
	assert.Equal(t, nil, err)
	err = bin.Delete()
	assert.Equal(t, nil, err)
}

func TestRoleBinding_Create(t *testing.T) {
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
				Kind : "User",
			}, {
				Kind : "User",
				Name: "shi",
			},
		},
	}
	role, err := bin.Create(rolebinding)
	log.Debug(role)
	log.Debug(err)
}

func TestRoleBinding_Update(t *testing.T) {
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
				Kind : "User",
			}, {
				Kind : "User",
				Name: "shi",
			},
		},
	}
	role, err := bin.Update(rolebinding)
	log.Debug(role)
	log.Debug(err)
}


package openshift

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/openshift/api/authorization/v1"
	"github.com/hidevopsio/hiboot/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPolicy_Create(t *testing.T) {
	name := ""
	namespace := ""
	policy, err := NewPolicy(name, namespace)
	assert.Equal(t, err, nil)
	po := &v1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	result, err := policy.Create(po)
	assert.Equal(t, err, nil)
	log.Debug(result)
}

func TestPolicy_Get(t *testing.T) {
	name := "admin"
	namespace := "demo-dev"
	policy, err := NewPolicy(name, namespace)
	assert.Equal(t, err, nil)
	po, err := policy.Get()
	assert.Equal(t, nil, err)
	log.Debug(po)
}
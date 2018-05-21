package k8s

import "testing"
import (
	core "k8s.io/api/core/v1"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConfigMapsCreate(t *testing.T) {
	name := "test"
	namespace := "demo-dev"
	data := map[string]string{}
	config := NewConfigMaps(name, namespace, data)
	config.Data = map[string]string{
	}
	result, err := config.Create()
	assert.Equal(t, nil, err)
	log.Info("", result)
}

func TestConfigMapsGet(t *testing.T) {
	name := "test1"
	namespace := "demo-dev"
	data := map[string]string{}
	config := NewConfigMaps(name, namespace, data)
	result, err := config.Get()
	assert.Equal(t, nil, err)
	log.Info(result)
}

func TestConfigMapsDelete(t *testing.T) {
	name := "test"
	namespace := "demo-dev"
	data := map[string]string{}
	config := NewConfigMaps(name, namespace, data)
	err := config.Delete()
	assert.Equal(t, nil, err)
}

func TestConfigMapsUpdate(t *testing.T) {
	name := "test"
	namespace := "demo-dev"
	data := map[string]string{}
	config := NewConfigMaps(name, namespace, data)
	configMap := &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			"default":"{a}",
		},
	}
	result, err := config.Update(configMap)
	assert.Equal(t, nil, err)
	log.Info("", result)
}
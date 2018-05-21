package k8s

import (
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"github.com/hidevopsio/hicicd/pkg/orch"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type ConfigMaps struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Data      map[string]string `json:"data"`
	Interface v1.ConfigMapInterface
}

func NewConfigMaps(name, namespace string, data map[string]string) *ConfigMaps {
	return &ConfigMaps{
		Name:      name,
		Namespace: namespace,
		Data:      data,
		Interface: orch.ClientSet.CoreV1().ConfigMaps(namespace),
	}
}

func (c *ConfigMaps) Create() (*core_v1.ConfigMap, error) {
	log.Debug("config map create :", c)
	configMap := &core_v1.ConfigMap{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: c.Name,
		},
		Data: c.Data,
	}
	co, err := c.Get()
	log.Debug("config map get :", co)
	if err == nil {
		con, err := c.Update(configMap)
		return con, err
	}
	config, er := c.Interface.Create(configMap)
	if er != nil {
		return nil, er
	}

	return config, nil
}

func (c *ConfigMaps) Get() (config *core_v1.ConfigMap, err error) {
	log.Info("get config map :", c.Name)
	result, err := c.Interface.Get(c.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *ConfigMaps) Delete() error {
	log.Info("get config map :", c.Name)
	err := c.Interface.Delete(c.Name, &meta_v1.DeleteOptions{})
	return err
}

func (c *ConfigMaps) Update(configMap *core_v1.ConfigMap) (*core_v1.ConfigMap, error) {
	log.Info("get config map :", c.Name)
	result, err := c.Interface.Update(configMap)
	return result, err
}

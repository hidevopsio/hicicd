// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package istio

import (
	"istio.io/istio/pilot/pkg/kube/inject"
	"istio.io/istio/pilot/pkg/model"
	"github.com/ghodss/yaml"
	"github.com/openshift/api/apps/v1"
	"fmt"
	"istio.io/istio/pilot/pkg/serviceregistry/kube"
	meshconfig "istio.io/api/mesh/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/orch"
	"io/ioutil"
)

type Injector struct {
	Version             string `json:"version"`
	Namespace           string `json:"namespace"`
	Hub                 string `json:"hub"`
	Tag                 string `json:"tag"`
	MeshConfigFile      string `json:"mesh_config_file"`
	InjectConfigFile    string `json:"inject_config_file"`
	MeshConfigMapName   string `json:"mesh_config_map_name"`
	InjectConfigMapName string `json:"inject_config_map_name"`
	DebugMode           bool   `json:"debug_mode"`
	SidecarProxyUID     uint64 `json:"sidecar_proxy_uid"`
	Verbosity           int    `json:"verbosity"`
	EnableCoreDump      bool   `json:"enable_core_dump"`
	ImagePullPolicy     string `json:"image_pull_policy"`
	IncludeIPRanges     string `json:"includeIPRanges"`
	ExcludeIPRanges     string `json:"exclude_ip_ranges"`
	IncludeInboundPorts string `json:"include_inbound_ports"`
	ExcludeInboundPorts string `json:"exclude_inbound_ports"`
	SidecarTemplate     string `json:"sidecar_template"`
}

const (
	configMapKey       = "mesh"
	injectConfigMapKey = "config"
)

//
//var (
//	hub                 string
//	tag                 string
//	sidecarProxyUID     uint64
//	verbosity           int
//	versionStr          string // override build version
//	enableCoreDump      bool
//	imagePullPolicy     string
//	includeIPRanges     string
//	excludeIPRanges     string
//	includeInboundPorts string
//	excludeInboundPorts string
//	debugMode           bool
//)


func (i *Injector) getMeshConfigFromConfigMap(kubeconfig string) (*meshconfig.MeshConfig, error) {
	_, client, err := kube.CreateInterface(kubeconfig)
	if err != nil {
		return nil, err
	}
	config, err := client.CoreV1().ConfigMaps(i.Namespace).Get(i.MeshConfigMapName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not read valid configmap %q from namespace  %q: %v - "+
			"Use --meshConfigFile or re-run kube-inject with `-i <istioSystemNamespace> and ensure valid MeshConfig exists",
			i.MeshConfigMapName, i.Namespace, err)
	}
	// values in the data are strings, while proto might use a
	// different data type.  therefore, we have to get a value by a
	// key
	configYaml, exists := config.Data[configMapKey]
	if !exists {
		return nil, fmt.Errorf("missing configuration map key %q", configMapKey)
	}
	return model.ApplyMeshConfigDefaults(configYaml)
}

func (i *Injector) getInjectConfigFromConfigMap(kubeconfig string) (string, error) {
	_, client, err := kube.CreateInterface(kubeconfig)
	if err != nil {
		return "", err
	}
	config, err := client.CoreV1().ConfigMaps(i.Namespace).Get(i.InjectConfigMapName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("could not find valid configmap %q from namespace  %q: %v - "+
			"Use --injectConfigFile or re-run kube-inject with `-i <istioSystemNamespace> and ensure istio-inject configmap exists",
			i.InjectConfigMapName, i.Namespace, err)
	}
	// values in the data are strings, while proto might use a
	// different data type.  therefore, we have to get a value by a
	// key
	injectData, exists := config.Data[injectConfigMapKey]
	if !exists {
		return "", fmt.Errorf("missing configuration map key %q in %q",
			injectConfigMapKey, i.InjectConfigMapName)
	}
	var injectConfig inject.Config
	if err := yaml.Unmarshal([]byte(injectData), &injectConfig); err != nil {
		return "", fmt.Errorf("unable to convert data from configmap %q: %v",
			i.InjectConfigMapName, err)
	}
	log.Debugf("using inject template from configmap %q", i.InjectConfigMapName)
	return injectConfig.Template, nil
}

func (i *Injector) Inject(in interface{}) (interface{}, error)  {

	mesh, err := i.getMeshConfigFromConfigMap(*orch.Kubeconfig)
	if err != nil {
		return nil, err
	}

	// get sidecar template from user's input, or get it from config map if user input is empty
	var sidecarTemplate string
	if i.SidecarTemplate != "" {
		sidecarTemplate = i.SidecarTemplate
	} else if i.InjectConfigFile != "" {
		injectionConfig, err := ioutil.ReadFile(i.InjectConfigFile) // nolint: vetshadow
		if err != nil {
			return nil, err
		}
		var config inject.Config
		if err := yaml.Unmarshal(injectionConfig, &config); err != nil {
			return nil, err
		}
		sidecarTemplate = config.Template
	} else if i.InjectConfigMapName != "" {
		if sidecarTemplate, err = i.getInjectConfigFromConfigMap(*orch.Kubeconfig); err != nil {
			return nil, err
		}
	} else {
		if sidecarTemplate, err = inject.GenerateTemplateFromParams(&inject.Params{
			InitImage:           inject.InitImageName(i.Hub, i.Tag, i.DebugMode),
			ProxyImage:          inject.ProxyImageName(i.Hub, i.Tag, i.DebugMode),
			Verbosity:           i.Verbosity,
			SidecarProxyUID:     i.SidecarProxyUID,
			Version:             i.Version,
			EnableCoreDump:      i.EnableCoreDump,
			Mesh:                mesh,
			ImagePullPolicy:     i.ImagePullPolicy,
			IncludeIPRanges:     i.IncludeIPRanges,
			ExcludeIPRanges:     i.ExcludeIPRanges,
			IncludeInboundPorts: i.IncludeInboundPorts,
			ExcludeInboundPorts: i.ExcludeInboundPorts,
			DebugMode:           i.DebugMode,
		}); err != nil {
			return nil, err
		}
	}
	log.Print(sidecarTemplate)

	out, err := inject.IntoObject(sidecarTemplate, mesh, in.(*v1.DeploymentConfig))
	return out, err
}
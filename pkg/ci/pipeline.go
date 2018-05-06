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

package ci

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/hidevopsio/hicicd/pkg/orch/k8s"
	"github.com/hidevopsio/hicicd/pkg/orch/openshift"
	"github.com/hidevopsio/hicicd/pkg/orch"
	"github.com/hidevopsio/hicicd/pkg/orch/istio"
	"github.com/jinzhu/copier"
	"path/filepath"
)

type Scm struct {
	Type string `json:"type"`
	Url  string `json:"url"`
	Ref  string `json:"ref"`
}

type DeploymentConfigs struct {
	HealthEndPoint string       `json:"health_end_point"`
	Skip           bool         `json:"skip"`
	ForceUpdate    bool         `json:"force_update"`
	Replicas       int32        `json:"replicas"`
	Env            []system.Env `json:"env"`
}

type BuildConfigs struct {
	Skip        bool         `json:"skip"` // TODO: ? Always, IfNotPresent, Never
	Tag         string       `json:"tag"`
	ImageStream string       `json:"image_stream"`
	Env         []system.Env `json:"env"`
}

type IstioConfigs struct {
	Skip                bool   `json:"skip"`
	Version             string `json:"version"`
	Namespace           string `json:"namespace"`
	DockerHub           string `json:"docker_hub"`
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
}

type Pipeline struct {
	Name              string            `json:"name" validate:"required"`
	App               string            `json:"app" validate:"required"`
	Profile           string            `json:"profile" validate:"required"`
	Project           string            `json:"project" validate:"required"`
	Namespace         string            `json:"namespace"`
	Scm               Scm               `json:"scm"`
	Version           string            `json:"version"`
	DockerRegistry    string            `json:"docker_registry"`
	Identifiers       []string          `json:"identifiers"`
	ConfigFiles       []string          `json:"config_files"`
	Ports             []orch.Ports      `json:"ports"`
	BuildConfigs      BuildConfigs      `json:"build_configs"`
	DeploymentConfigs DeploymentConfigs `json:"deployment_configs"`
	IstioConfigs      IstioConfigs      `json:"istio_configs"`
}

type Configuration struct {
	Pipeline Pipeline `mapstructure:"pipeline"`
}

// @Title Init
// @Description set default value
// @Param pipeline
// @Return error
func (p *Pipeline) Init(pl *Pipeline) {
	log.Debug("Pipeline.EnsureParam()")

	// load config file
	if pl != nil {
		/*		builder := &Builder{}
				c := builder.Build(pl.Name)*/

		b := &system.Builder{
			Path:       filepath.Join(utils.GetWorkDir(), "config"),
			Name:       "pipeline",
			FileType:   "yaml",
			Profile:    pl.Name,
			ConfigType: Configuration{},
		}
		cp, err := b.Build()
		if err != nil {
			return
		}
		c := cp.(*Configuration)
		if pl.Profile == "" {
			p.Profile = "dev"
		}

		mergo.Merge(&c.Pipeline, pl, mergo.WithOverride)
		mergo.Merge(p, c.Pipeline, mergo.WithOverride)

	}

	utils.Replace(p, p)

	if "" == p.Namespace {
		if "" == pl.Profile {
			p.Namespace = p.Project
		} else {
			p.Namespace = p.Project + "-" + p.Profile
		}
	}

	log.Debug(p)
}

func (p *Pipeline) CreateSecret(username, password string, isToken bool) (string, error) {
	log.Debug("Pipeline.CreateSecret()")
	if username == "" {
		return "", fmt.Errorf("unkown username")
	}
	// Create secret
	secretName := username + "-secret"
	secret := k8s.NewSecret(secretName, username, password, p.Namespace, isToken)
	err := secret.Create()

	return secretName, err
}

func (p *Pipeline) Build(secret string, completedHandler func() error) error {
	log.Debug("Pipeline.Build()")

	if p.BuildConfigs.Skip {
		return completedHandler()
	}

	scmUrl := p.CombineScmUrl()
	buildConfig, err := openshift.NewBuildConfig(p.Namespace, p.App, scmUrl, p.Scm.Ref, secret, p.Version, p.BuildConfigs.ImageStream)
	if err != nil {
		return err
	}
	_, err = buildConfig.Create()
	if err != nil {
		return err
	}
	// Build image stream
	build, err := buildConfig.Build(p.BuildConfigs.Env)

	if err != nil {
		return err
	}

	buildConfig.Watch(build, completedHandler)

	return err
}

func (p *Pipeline) CombineScmUrl() string {
	scmUrl := p.Scm.Url + "/" + p.Project + "/" + p.App + "." + p.Scm.Type
	return scmUrl
}

func (p *Pipeline) RunUnitTest() error {
	log.Debug("Pipeline.RunUnitTest()")
	return nil
}

func (p *Pipeline) RunIntegrationTest() error {
	log.Debug("Pipeline.RunIntegrationTest()")
	return nil
}

func (p *Pipeline) Analysis() error {
	log.Debug("Pipeline.Analysis()")
	return nil
}

func (p *Pipeline) CreateDeploymentConfig(force bool, injectSidecar func(in interface{}) (interface{}, error)) error {
	log.Debug("Pipeline.CreateDeploymentConfig()")

	// new dc instance
	dc, err := openshift.NewDeploymentConfig(p.App, p.Namespace, p.Version)
	if err != nil {
		return err
	}

	err = dc.Create(&p.DeploymentConfigs.Env, &p.Ports, p.DeploymentConfigs.Replicas, force, p.DeploymentConfigs.HealthEndPoint, injectSidecar)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pipeline) Deploy() error {
	log.Debug("Pipeline.Deploy()")

	// new dc instance
	dc, err := openshift.NewDeploymentConfig(p.App, p.Namespace, p.Version)
	if err != nil {
		return err
	}

	d, err := dc.Instantiate()
	log.Debug(d.Name)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pipeline) CreateService() error {
	log.Debug("Pipeline.CreateService()")

	// new dc instance
	svc := k8s.NewService(p.App, p.Namespace)

	err := svc.Create(&p.Ports)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pipeline) CreateRoute() error {
	log.Debug("Pipeline.CreateRoute()")

	route, err := openshift.NewRoute(p.App, p.Namespace)
	if err != nil {
		return err
	}

	err = route.Create(8080)
	return nil
}

func (p *Pipeline) Run(username, password string, isToken bool) error {
	log.Debug("Pipeline.Run()")
	// TODO: first, let's check if namespace is exist or not

	// TODO: check if the same app in the same namespace is already in running status.

	// create secret for building image
	secret, err := p.CreateSecret(username, password, isToken)
	if err != nil {
		return fmt.Errorf("failed on CreateSecret! %s", err.Error())
	}

	// build image
	err = p.Build(secret, func() error {

		if !p.DeploymentConfigs.Skip {

			// create dc - deployment config
			err = p.CreateDeploymentConfig(p.DeploymentConfigs.ForceUpdate, func(in interface{}) (interface{}, error) {
				if p.IstioConfigs.Skip {
					return in, nil
				}
				
				injector := &istio.Injector{}
				copier.Copy(injector, p.IstioConfigs)
				return injector.Inject(in)
			})
			if err != nil {
				log.Error(err.Error())
				return fmt.Errorf("failed on CreateDeploymentConfig! %s", err.Error())
			}

			//// deploy
			//err = p.Deploy()
			//if err != nil {
			//	log.Error(err.Error())
			//	return fmt.Errorf("failed on Deploy! %s", err.Error())
			//}

			rc := k8s.NewReplicationController(p.App, p.Namespace)
			// rc.Watch(message, handler)
			err = rc.Watch(func() error {
				log.Debug("Completed!")
				return nil
			})
			if err != nil {
				log.Error(err.Error())
				return fmt.Errorf("failed on watch rc! %s", err.Error())
			}
		}

		// create service
		err = p.CreateService()
		if err != nil {
			log.Error(err.Error())
			return fmt.Errorf("failed on CreateService! %s", err.Error())
		}

		// create route
		err = p.CreateRoute()
		if err != nil {
			log.Error(err.Error())
			return fmt.Errorf("failed on CreateRoute! %s", err.Error())
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed on Build! %s", err.Error())
	}

	// interact with developer
	// deploy to test ? yes/no

	// finally, all steps are done well, let tell the client ...
	return nil
}

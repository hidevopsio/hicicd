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
	authorization_v1 "github.com/openshift/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"github.com/hidevopsio/hicicd/pkg/auth"
	"os"
	"strings"
	"encoding/json"
	"github.com/kevholditch/gokong"
)

type Scm struct {
	Type string `json:"type"`
	Url  string `json:"url"`
	Ref  string `json:"ref"`
}

type DeploymentConfigs struct {
	HealthEndPoint string       `json:"health_end_point"`
	Enable         bool         `json:"enable"`
	ForceUpdate    bool         `json:"force_update"`
	Replicas       int32        `json:"replicas"`
	Env            []system.Env `json:"env"`
	Labels         Labels       `json:"labels"`
	Project        string       `json:"project"`
}

type Labels struct {
	App     string `json:"app"`
	Version string `json:"version"`
	Cluster string `json:"cluster"`
}

type BuildConfigs struct {
	Enable      bool         `json:"enable"` // TODO: ? Always, IfNotPresent, Never
	TagFrom     string       `json:"tag_from"`
	ImageStream string       `json:"image_stream"`
	Env         []system.Env `json:"env"`
	Rebuild     bool         `json:"rebuild"`
	Project     string       `json:"project"`
	Namespace   string       `json:"namespace"`
	Branch      string       `json:"branch"`
}

type IstioConfigs struct {
	Enable              bool   `json:"enable"`
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
	Profile           string            `json:"profile"`
	Project           string            `json:"project" validate:"required"`
	Cluster           string            `json:"cluster"`
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
	GatewayConfigs    GatewayConfigs    `json:"gateway_configs"`
}

type Configuration struct {
	Pipeline Pipeline `mapstructure:"pipeline"`
}

type GatewayConfigs struct {
	Enable      bool   `json:"enable"`
	Uri         string `json:"uri"`
	UpstreamUrl string `json:"upstream_url"`
}

// @Title Init
// @Description set default value
// @Param pipeline
// @Return error
func (p *Pipeline) Init(pl *Pipeline) {
	log.Debug("Pipeline.EnsureParam()")
	// load config file
	if pl != nil {

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

	if "" == p.BuildConfigs.Namespace {
		if "" == pl.Profile {
			p.BuildConfigs.Namespace = p.BuildConfigs.Project
		} else {
			p.BuildConfigs.Namespace = p.BuildConfigs.Project + "-" + p.Profile
		}
	}

	if "" == p.Profile {
		p.BuildConfigs.TagFrom = ""
	} else if p.Profile == "dev" || p.Profile == "test" {
		p.BuildConfigs.TagFrom = "dev"
	} else {
		p.BuildConfigs.TagFrom = "stage"
	}
	if pl.Scm.Ref == "master" {
		p.BuildConfigs.Branch = ""
	} else {
		p.BuildConfigs.Branch = pl.Scm.Ref
	}
	p.GatewayConfigs.UpstreamUrl = p.App + "." + p.Namespace + ":8080"
	if pl.DeploymentConfigs.Enable == false {
		p.DeploymentConfigs.Enable = false
	}
	if pl.BuildConfigs.Enable == false {
		p.BuildConfigs.Enable = false
	}
	if pl.GatewayConfigs.Enable == false {
		p.GatewayConfigs.Enable = false
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
	secret := k8s.NewSecret(secretName, username, password, p.BuildConfigs.Namespace, isToken)
	err := secret.Create()

	return secretName, err
}

func (p *Pipeline) CreateProject() error {
	project, err := openshift.NewProject(p.Namespace, "", "")
	if err != nil {
		log.Error("create namespace err :", err)
		return err
	}
	prj, err := project.Get()
	if err == nil {
		log.Debug("exists project debug:", prj.Name)
		return err
	}
	newProject, err := project.Create()
	if err != nil {
		return err
	}
	//init namespace
	err = p.InitProject()
	if err != nil {
		return err
	}
	log.Debug("create project debug", newProject)
	return nil
}

func (p *Pipeline) CreateRoleBinding(username, metaName, roleRefName string) error {
	roleBinding, err := openshift.NewRoleBinding(metaName, p.Namespace)
	if err != nil {
		return err
	}
	r, err := roleBinding.Get()
	if err != nil {
		role := &authorization_v1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      metaName,
				Namespace: p.Namespace,
			},
			RoleRef: corev1.ObjectReference{
				Name: roleRefName,
			},
			Subjects: []corev1.ObjectReference{
				{
					Name: username,
					Kind: "User",
				},
			},
		}
		_, err = roleBinding.Create(role)
		return err
	}
	for _, value := range r.Subjects {
		if value.Name == username {
			return nil
		}
	}
	reference := corev1.ObjectReference{
		Name: username,
		Kind: "User",
	}
	r.Subjects = append(r.Subjects, reference)
	r.UserNames = append(r.UserNames, username)
	_, err = roleBinding.Create(r)
	return err
}

func (p *Pipeline) Build(secret string, completedHandler func() error) error {
	log.Debug("Pipeline.Build()")

	if !p.BuildConfigs.Enable {
		return completedHandler()
	}

	scmUrl := p.CombineScmUrl()
	buildConfig, err := openshift.NewBuildConfig(p.BuildConfigs.Namespace, p.App, scmUrl, p.BuildConfigs.Branch, secret, p.Version, p.BuildConfigs.ImageStream, p.BuildConfigs.Rebuild)
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
		log.Error("buildConfig.Build(p.BuildConfigs.Env)", err)
		return err
	}

	err = buildConfig.Watch(build, completedHandler)
	log.Info("buildConfig.Watch", err)
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
	var l map[string]string
	labels, _ := json.Marshal(p.DeploymentConfigs.Labels)
	err = json.Unmarshal(labels, &l)
	err = dc.Create(&p.DeploymentConfigs.Env, l, &p.Ports, p.DeploymentConfigs.Replicas, force, p.DeploymentConfigs.HealthEndPoint, p.Profile, injectSidecar)
	if err != nil {
		log.Error("dc.Create ", err)
		return err
	}

	return nil
}

func (p *Pipeline) InstantiateDeploymentConfig() error {
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

func (p *Pipeline) CreateKongGateway(upstreamUrl string) error {
	log.Debug("Pipeline.CreateKongGateway()")
	uris := p.GatewayConfigs.Uri
	host := os.Getenv("KONG_HOST")
	host = strings.Replace(host, "${profile}", p.Profile, -1)
	hosts := strings.Split(host, ",")
	apiRequest := &gokong.ApiRequest{
		Name:                   p.App + "-" + p.Project,
		Hosts:                  hosts,
		Uris:                   []string{uris},
		UpstreamUrl:            "http://" + upstreamUrl,
		StripUri:               false,
		PreserveHost:           true,
		Retries:                "3",
		UpstreamConnectTimeout: 1000,
		UpstreamSendTimeout:    2000,
		UpstreamReadTimeout:    3000,
		HttpsOnly:              false,
		HttpIfTerminated:       true,
	}
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	baseUrl = strings.Replace(baseUrl, "${profile}", p.Profile, -1)
	config := &gokong.Config{
		HostAddress: baseUrl,
	}
	_, err := gokong.NewClient(config).Apis().Create(apiRequest)
	return err
}

func (p *Pipeline) CreateService() error {
	log.Debug("Pipeline.CreateService()")

	// new dc instance
	svc := k8s.NewService(p.App, p.Namespace)

	err := svc.Create(&p.Ports)
	log.Debug("Pipeline.CreateService svc.Create result: ", err)
	return err
}

func (p *Pipeline) CreateRoute() (string, error) {
	log.Debug("Pipeline.CreateRoute()")
	upstreamUrl := ""
	route, err := openshift.NewRoute(p.App, p.Namespace)
	if err != nil {
		return upstreamUrl, err
	}

	upstreamUrl, err = route.Create(8080)
	return upstreamUrl, err
}

func (p *Pipeline) GetImageStreamTag() error {
	log.Debug("pipeline get image stream tag :")
	ist, err := openshift.NewImageStreamTags(p.App, p.Version, p.BuildConfigs.Project+"-"+p.BuildConfigs.TagFrom)
	if err != nil {
		log.Error("Pipeline.CreateImageStreamTag.NewImageStreamTags", err)
		return err
	}
	image, err := ist.Get()
	if err != nil {
		log.Error("images not found :", err)
		return err
	}
	log.Debug("pipeline get image stream tag name:", image.Name)
	return nil
}

func (p *Pipeline) CreateImageStreamTag() error {
	log.Debug("Pipeline.CreateImageStreamTag")
	ist, err := openshift.NewImageStreamTags(p.App, p.Version, p.Namespace)
	if err != nil {
		log.Error("Pipeline.CreateImageStreamTag.NewImageStreamTags", err)
		return err
	}
	_, err = ist.Create(p.BuildConfigs.Project + "-" + p.BuildConfigs.TagFrom)
	return err
}

func (p *Pipeline) InitProject() error {
	//init Groups
	roleBinding := &openshift.RoleBinding{
		Namespace: p.Namespace,
	}
	//init image builders
	err := roleBinding.InitImageBuilders()
	if err != nil {
		return err
	}
	//init image pullers
	err = roleBinding.InitImagePullers()
	if err != nil {
		return err
	}
	//init system deployers
	err = roleBinding.InitSystemDeployers()
	if err != nil {
		return err
	}
	return nil
}

func (p *Pipeline) Deploy() error {
	log.Info("p.Deploy()")
	if p.DeploymentConfigs.Enable {

		// create dc - deployment config
		err := p.CreateDeploymentConfig(p.DeploymentConfigs.ForceUpdate, func(in interface{}) (interface{}, error) {
			if !p.IstioConfigs.Enable {
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
	err := p.CreateService()
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("failed on CreateService! %s", err.Error())
	}

	// create route
	upstreamUrl, err := p.CreateRoute()
	log.Debug("CreateRoute get upstream url :", upstreamUrl)
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("failed on CreateRoute! %s", err.Error())
	}

	//create kong-gateway
	if p.GatewayConfigs.Enable {
		err = p.CreateKongGateway(p.GatewayConfigs.UpstreamUrl)
	}
	return err
}

func (p *Pipeline) Run(username, password, token string, uid int, isToken bool) error {
	log.Debug("Pipeline.Run()")
	// TODO: check if the same app in the same namespace is already in running status.
	permission := &auth.Permission{}
	metaName, roleRefName, accessLevelValue, err := permission.Get(p.Scm.Url, token, p.App, p.BuildConfigs.Project, uid)
	// TODO: accessLevelValue permission 30
	if err != nil || accessLevelValue < auth.DeveloperPermissions {
		return err
	}
	// TODO: first, let's check if namespace is exist or not
	//create  namespace
	err = p.CreateProject()
	if err != nil {
		log.Error("Pipeline run new namespace err:", err)
		return err
	}
	//Authorize the user
	err = p.CreateRoleBinding(username, metaName, roleRefName)
	if err != nil {
		log.Error("Pipeline run Create RoleBinding err :", err)
		return err
	}
	if p.BuildConfigs.TagFrom == p.Profile && p.BuildConfigs.Project == p.Project {

		// create secret for building image
		secret, err := p.CreateSecret(username, password, isToken)
		if err != nil {
			return fmt.Errorf("failed on CreateSecret! %s", err.Error())
		}

		// build image
		err = p.Build(secret, func() error {
			return p.Deploy()
		})
	} else {
		//TODO check if tag from is exist or not
		err = p.GetImageStreamTag()
		if err != nil {
			return nil
		}
		err = p.CreateImageStreamTag()
		if err != nil {
			log.Error("Pipeline.Run.CreateImageStreamTag error", err)
			return err
		}
		p.Deploy()
	}

	if err != nil {
		return fmt.Errorf("failed on Build! %s", err.Error())
	}

	// interact with developer
	// deploy to test ? yes/no

	// finally, all steps are done well, let tell the client ...
	return nil
}

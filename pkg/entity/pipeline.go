package entity

import (
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hicicd/pkg/orch"
	"path/filepath"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/imdario/mergo"
	"github.com/hidevopsio/hiboot/pkg/log"
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
	Enable      bool         `json:"enable"`
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
	NodeSelector      string            `json:"node_selector"`
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
func (p *Pipeline) Init(pl *Pipeline, selector *Selector) {
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
	for _, node := range selector.Nodes {
		if node.Profile == p.Profile {
			p.NodeSelector = node.NodeSelector
		}
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
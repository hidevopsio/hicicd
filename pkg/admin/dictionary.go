package admin

type Dictionary struct {
	Id                string            `json:"id"`
	Scm               Scm               `json:"scm"`
	Istio             Istio             `json:"istio"`
	BuildConfigs      BuildConfigs      `json:"build_configs"`
	DeploymentConfigs DeploymentConfigs `json:"deployment_configs"`
	Profiles          []string          `json:"profiles"`
	ImageStreamTags   ImageStreamTags   `json:"image_stream_tags"`
	Version           string            `json:"version"`
	Url               string            `json:"url"`
}

type ImageStreamTags struct {
	NodeImageStreamTags []string `json:"node_selector"`
	JavaSelector        []string `json:"java_selector"`
}

type DeploymentConfigs struct {
	Enable      bool `json:"enable"`
	ForceUpdate bool `json:"force_update"`
}

type BuildConfigs struct {
	Enable bool `json:"enable"`
}

type Istio struct {
	Enable bool `json:"enable"`
}

type Scm struct {
	Branches []string `json:"branches"`
}

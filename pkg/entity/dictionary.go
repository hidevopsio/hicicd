package entity

type Dictionary struct {
	Id                string             `json:"id"`
	Scm               scm                `json:"scm"`
	Istio             Istio              `json:"istio"`
	BuildConfigs      buildConfigs       `json:"build_configs"`
	DeploymentConfigs deploymentConfigs  `json:"deployment_configs"`
	Profiles          []string           `json:"profiles"`
	ImageStreamTags   map[string][]Image `json:"image_stream_tags"`
	Version           string             `json:"version"`
	Hosts             []host             `json:"hosts"`
}

type host struct {
	Profile string `json:"profile"`
	Host    string `json:"host"`
}

type Image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Text       string `json:"text"`
	Value      string `json:"value"`
	Name       string `json:"name"`
}

type deploymentConfigs struct {
	Enable      bool `json:"enable"`
	ForceUpdate bool `json:"force_update"`
}

type buildConfigs struct {
	Enable bool `json:"enable"`
}

type Istio struct {
	Enable bool `json:"enable"`
}

type scm struct {
	Branches []string `json:"branches"`
}

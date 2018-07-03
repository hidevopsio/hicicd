package admin

type Dictionary struct {
	Id                string             `json:"id"`
	Scm               Scm                `json:"scm"`
	Istio             Istio              `json:"istio"`
	BuildConfigs      BuildConfigs       `json:"build_configs"`
	DeploymentConfigs DeploymentConfigs  `json:"deployment_configs"`
	Profiles          []string           `json:"profiles"`
	ImageStreamTags   map[string][]Image `json:"image_stream_tags"`
	Version           string             `json:"version"`
	Hosts             []Host             `json:"hosts"`
}

type Host struct {
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

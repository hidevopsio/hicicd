package impl

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/log"
	"os"
	"github.com/hidevopsio/hicicd/pkg/ci"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hicicd/pkg/scm/gitlab"
)

func TestNodeJsPipeline(t *testing.T) {

	log.Debug("Test NodeJs Pipeline")

	nodeJs := &NodeJsPipeline{}
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	pi := &ci.Pipeline{
		Name:    "nodejs",
		Project: "demo",
		Profile: "dev",
		App:     "hello-angular",
		Scm:     ci.Scm{Url: os.Getenv("SCM_URL")},
		DeploymentConfigs: ci.DeploymentConfigs{
			ForceUpdate: true,
		},
		//DeploymentConfigs: ci.DeploymentConfigs{
		//	Skip: true,
		//},
		//BuildConfigs: ci.BuildConfigs{
		//	Skip: true,
		//},
	}
	nodeJs.Init(pi)
	baseUrl :=  os.Getenv("SCM_URL")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	log.Debug(gs)
	assert.Equal(t, username, gs.Username)
	err = nodeJs.Run(username, password, gs.PrivateToken, gs.ID, false)
	assert.Equal(t, nil, err)
}

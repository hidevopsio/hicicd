package gitlab

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/magiconair/properties/assert"
	"github.com/xanzy/go-gitlab"
)

func TestGetProjectMember(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	id := 1024
	projectMember := &ProjectMember{
		Token:   gs.PrivateToken,
		BaseUrl: baseUrl,
		Pid:     id,
		User:    401,
	}
	p, err := projectMember.GetProjectMember()
	assert.Equal(t, p.AccessLevel, gitlab.MasterPermissions)
	log.Info(p.AccessLevel)

}

package gitlab

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/magiconair/properties/assert"
	"github.com/xanzy/go-gitlab"
	"net/http"
)

func TestGetProjectMember(t *testing.T) {
	baseUrl := "http://gitlab.vpclub.cn:8022"
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	c := gitlab.NewClient(&http.Client{}, gs.PrivateToken)
	c.SetBaseURL(baseUrl + ApiVersion)
	groupMember, _, err := c.GroupMembers.GetGroupMember("demo",gs.ID)
	log.Debug(err)
	log.Debug(groupMember)
	//p, err := projectMember.GetProjectMember(gs.PrivateToken, baseUrl, id, gs.ID, gid)
}

func TestListProjectMembers(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	id := 1024
	projectMember := new(ProjectMember)
	pm, err := projectMember.ListProjectMembers(gs.PrivateToken, baseUrl, id)
	assert.Equal(t, nil, err)
	log.Info(pm)
}
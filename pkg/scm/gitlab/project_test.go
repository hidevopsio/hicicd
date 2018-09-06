package gitlab

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/magiconair/properties/assert"
)

func TestGetProject(t *testing.T) {
	baseUrl := "http://gitlab.vpclub.cn:8022"
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	project := new(Project)
	id := "demo/hello-world"
	pid, gid, err := project.GetProject(baseUrl, id, gs.PrivateToken)
	assert.Equal(t, nil, err)
	assert.Equal(t, 905, pid)
	assert.Equal(t, 4, gid)
}

func TestListProjects(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	page := 1
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	project := new(Project)
	projects, err := project.ListProjects(baseUrl, gs.PrivateToken, "demo",page)
	assert.Equal(t, nil, err)
	assert.Equal(t, username, gs.Username)
	log.Debug("project size:", len(projects))
	log.Debug("project size:", projects)
}

func TestSearch(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	query := "hello-world"
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	project := new(Project)
	p, err := project.Search(baseUrl, gs.PrivateToken, query)
	assert.Equal(t, nil, err)
	assert.Equal(t, username, gs.Username)
	log.Debug(p)
}
package gitlab

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/magiconair/properties/assert"
)

func TestGetProject(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	id := 1067
	project := new(Project)
	p, err := project.GetProject()
	if err != nil {
		return
	}
	assert.Equal(t, nil, err)
	assert.Equal(t, p.ID, id)
}

func TestListProjects(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	project := new(Project)
	projects, err := project.ListProjects()
	assert.Equal(t, nil, err)
	assert.Equal(t, username, gs.Username)
	log.Debug("project size:", len(projects))
	log.Debug("project size:", projects)
}

func TestListUserProjects(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	name := "hello-world"
	namespace := "demo"
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	project := new(Project)
	assert.Equal(t, nil, err)
	pid, err := project.ListUserProjects(baseUrl, gs.PrivateToken, name, namespace)
	log.Debug("get user project pid :", pid)
	assert.Equal(t, true, pid)
	assert.Equal(t, nil, err)
}

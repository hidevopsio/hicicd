package gitlab

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/magiconair/properties/assert"
)

func TestGetProject(t *testing.T) {
	baseUrl :=  os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	id := 1067
	token := gs.PrivateToken
	project, err := gs.GetProject(token, baseUrl, id)
	if err != nil {
		return
	}
	assert.Equal(t, nil, err)
	assert.Equal(t, project.ID, id)
}

func TestGetProjectList(t *testing.T) {
	baseUrl :=  os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	projects, err := gs.GetProjectLlist(gs.PrivateToken, baseUrl)
	assert.Equal(t, nil, err)
	assert.Equal(t, username, gs.Username)
	log.Debug("project size:", len(projects))
	log.Debug("project size:", projects)
}

func TestGetUserProject(t *testing.T) {
	baseUrl :=  os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	name := "hello-world"
	namespace := "demo"
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	exists := gs.GetUserProject(gs.PrivateToken, baseUrl, name, namespace)
	log.Debug("get user project exists :", exists)
	assert.Equal(t, true, exists)
}
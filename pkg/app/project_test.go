package app

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm/gitlab"
	"github.com/magiconair/properties/assert"
)

func TestProjects(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	page   := 1
	search := "hello-world"
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	p := new(Project)
	_, err = p.ListProjects(gs.PrivateToken, baseUrl, search, page)
	assert.Equal(t, nil, err)
}



func TestProjectMember(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	pid := 905
	gid := 4
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	p := new(Project)
	_, err = p.GetProjectMember(gs.PrivateToken, baseUrl, pid, gs.ID, gid)
	assert.Equal(t, nil, err)
}
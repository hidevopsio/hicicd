package gitlab

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"os"
)

func TestListGroups(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	group := new(Group)
	g, err := group.ListGroups(gs.PrivateToken, baseUrl)
	assert.Equal(t, nil, err)
	log.Infof("groups :%v", g)
}

func TestGetGroup(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	gid := 158
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	group := new(Group)
	scmGroup, err := group.GetGroup(gs.PrivateToken, baseUrl, gid)
	assert.Equal(t, nil, err)
	log.Info(scmGroup)
}

func TestListGroupProjects1(t *testing.T)   {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	gid := 158
	page := 1
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	group := new(Group)
	scmProject, err := group.ListGroupProjects(gs.PrivateToken, baseUrl, gid, page)
	assert.Equal(t, nil, err)
	log.Info(scmProject)
}


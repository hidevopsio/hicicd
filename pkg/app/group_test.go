package app

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm/gitlab"
	"github.com/magiconair/properties/assert"
)

func TestListGroups(t *testing.T)  {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	page := 6
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	g := new(Group)
	group, err := g.ListGroups(gs.PrivateToken, baseUrl, gs.ID, page)
	assert.Equal(t, nil, err)
	log.Info(group)
}

func TestGetGroup(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	gid := 574
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	g := new(Group)
	group, err := g.GetGroup(gs.PrivateToken, baseUrl, gid)
	assert.Equal(t, nil, err)
	log.Info(group)
}

func TestGroupMembers(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	gid := 574
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	g := new(Group)
	_, err = g.ListGroupMembers(gs.PrivateToken, baseUrl, gid, gs.ID)
	assert.Equal(t, nil, err)
}

func TestGroupProjects(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	gid := 574
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	g := new(Group)
	_, err = g.ListGroupProjects(gs.PrivateToken, baseUrl, gid, gs.ID)
	assert.Equal(t, nil, err)
}
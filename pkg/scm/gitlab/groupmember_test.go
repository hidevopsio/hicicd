package gitlab

import (
"testing"
"os"
"github.com/hidevopsio/hiboot/pkg/log"
"github.com/magiconair/properties/assert"
)

func TestListGroupMembers(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	gid := 4
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	groupMember := GroupMember{}
	g, err := groupMember.ListGroupMembers(gs.PrivateToken, baseUrl, gid, gs.ID)
	assert.Equal(t, nil, err)
	log.Info(g)
}

func TestGetGroupMember(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	gid := 4
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	groupMember := new(GroupMember)
	gm, err := groupMember.GetGroupMember(gs.PrivateToken, baseUrl, gid, gs.ID)
	assert.Equal(t, nil, err)
	log.Info(gm)
}


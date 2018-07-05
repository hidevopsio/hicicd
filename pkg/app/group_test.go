package app

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm/gitlab"
)

func TestListGroups(t *testing.T)  {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	gid := 158
	page := 1
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	group := new(Group)
}
package gitlab

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/magiconair/properties/assert"
)

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
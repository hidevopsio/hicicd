package info

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm/gitlab"
	"github.com/magiconair/properties/assert"
	"fmt"
)


func TestRepositoryType(t *testing.T) {
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	t1 := new(TypeInfo)
	ref := "master"
	pid := 1151
	err = t1.RepositoryType(baseUrl, gs.PrivateToken, ref, pid)
	assert.Equal(t, nil, err)
	fmt.Printf("log : %v", t1)

}

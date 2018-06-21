package gitlab

import (
	"testing"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/magiconair/properties/assert"
)

func TestGetRepositoty(t *testing.T){
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	query := "pom.xml"
	ref := "master"
	pid := 905
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	repository := new(Repository)
	repository.GetRepository(baseUrl, gs.PrivateToken, query, ref, pid)
	assert.Equal(t, nil, err)
	assert.Equal(t, username, gs.Username)

}

func TestListTree(t *testing.T){
	baseUrl := os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	ref := "master"
	pid := 905
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	repository := new(Repository)
	tr , err := repository.ListTree(baseUrl, gs.PrivateToken, ref, pid)
	assert.Equal(t, nil, err)
	assert.Equal(t, username, gs.Username)
	log.Info(tr)
}
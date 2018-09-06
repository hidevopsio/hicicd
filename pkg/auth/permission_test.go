package auth

import (
	"testing"
	"os"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func TestGet(t *testing.T)  {
	baseUrl :=  os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	name := "hello-world"
	namespace := "demo"
	u := new(User)
	token, id, message, err := u.Login(baseUrl, username,  password)
	assert.Equal(t, nil, err)
	log.Debug(id)
	log.Debug(token)
	log.Debug(message)
	permission := new(Permission)
	projectMember, roleRefName,  accessLevelValue, err :=permission.Get(baseUrl, token, name, namespace, id)
	assert.Equal(t, nil, err)
	log.Debug("projectMember: ",projectMember, " roleRefName: ", roleRefName, " accessLevelValue: ", accessLevelValue)
}

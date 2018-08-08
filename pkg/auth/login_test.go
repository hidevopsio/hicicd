package auth

import (
	"testing"
	"os"
	"github.com/magiconair/properties/assert"
)

func TestGetUser(t *testing.T){
	baseUrl := os.Getenv("SCM_URL")
	accessToken := "0767a31c454f00c9a7682871732b4ddb94be611a87717ae12e3b8db2842c3146"
	l := new(Login)
	user, err := l.GetUser(baseUrl, accessToken)
	assert.Equal(t, nil, err)
	assert.Equal(t, "chulei", user.Username)
}

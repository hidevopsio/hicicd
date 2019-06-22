package controllers

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"os"
	"time"
	"fmt"
	"net/http"
	"github.com/magiconair/properties/assert"
)

func TestGetProject(t *testing.T) {
	app := web.NewTestApplication(t, new(ProjectController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	assert.Equal(t, nil, err)
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Get("/project").WithHeader("Authorization", bt).WithQuery("page", 1).Expect().Status(http.StatusOK)
}

func TestGetMember(t *testing.T) {
	app := web.NewTestApplication(t, new(ProjectController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	assert.Equal(t, nil, err)
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Get("/project/member").WithHeader("Authorization", bt).WithQuery("pid", 105).WithQuery("gid", 4).Expect().Status(http.StatusOK)
}
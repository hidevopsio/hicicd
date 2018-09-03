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

func TestGroup(t *testing.T) {
	app := web.NewTestApplication(t, new(GroupController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	assert.Equal(t, nil, err)
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Get("/group").WithHeader("Authorization", bt).WithQuery("page", 1).Expect().Status(http.StatusOK)
}

func TestProject(t *testing.T) {
	app := web.NewTestApplication(t, new(GroupController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	assert.Equal(t, nil, err)
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Get("/group/projects").WithHeader("Authorization", bt).WithQuery("page", 1).WithQuery("gid", 158).Expect().Status(http.StatusOK)
}
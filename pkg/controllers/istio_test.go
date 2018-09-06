package controllers

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"os"
	"time"
	"fmt"
	"net/http"
	"github.com/magiconair/properties/assert"
)

func TestGetIstio(t *testing.T) {
	app := web.NewTestApplication(t, new(IstioController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	assert.Equal(t, nil, err)
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Get("/istio").WithHeader("Authorization", bt).WithQuery("type", 1).
		WithQuery("name", "hello-world").
		WithQuery("namespace", "demo-stage").Expect().Status(http.StatusOK)
}

func TestDeleteIstio(t *testing.T) {
	app := web.NewTestApplication(t, new(IstioController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	assert.Equal(t, nil, err)
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Delete("/istio").WithHeader("Authorization", bt).WithQuery("type", 1).
		WithQuery("name", "hello-world").
		WithQuery("namespace", "demo-stage").Expect().Status(http.StatusOK)
}
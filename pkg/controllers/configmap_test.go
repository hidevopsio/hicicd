package controllers

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"time"
	"fmt"
	"github.com/magiconair/properties/assert"
	"os"
	"net/http"
	"github.com/hidevopsio/hioak/pkg/k8s"
)

func TestCreate(t *testing.T) {
	app := web.NewTestApplication(t, new(ConfigMapController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	assert.Equal(t, nil, err)
	cm := k8s.ConfigMaps{
		Name: "demo",
		Namespace: "demo-stage",
		Data: map[string]string{
			"a":"data",
		},
	}
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Post("/configMap").WithHeader("Authorization", bt).WithJSON(cm).Expect().Status(http.StatusOK)
}

func TestGet(t *testing.T)  {
	app := web.NewTestApplication(t, new(ConfigMapController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Minute)
	assert.Equal(t, nil, err)
	name := "demo"
	namespace := "demo-stage"
	//path := fmt.Sprintf("/configMap?name=%s&namespace=%s", url.QueryEscape(name), url.QueryEscape(namespace))
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Get("/configMap").WithHeader("Authorization", bt).WithQuery("name", name).WithQuery("namespace", namespace).Expect().Status(http.StatusOK)
}
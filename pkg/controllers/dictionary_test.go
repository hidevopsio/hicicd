package controllers

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"os"
	"time"
	"fmt"
	"net/http"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hicicd/pkg/entity"
)

func TestPostDictionary(t *testing.T) {
	app := web.NewTestApplication(t, new(DictionaryController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	assert.Equal(t, nil, err)
	dictionary := &entity.Dictionary{
		Id:"1",
		Profiles: []string{"dev"},
		Istio: entity.Istio{
			Enable:false,
		},
		Version:"v1",
		ImageStreamTags: map[string][]entity.Image{
			"java": []entity.Image{
				entity.Image{
					Text: "java",
					Repository: "s2i-java",
					Tag :"1.0.5",
					Value: "s2i-java:latest",
					Name: "java",
				},
			},
		},
	}
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Post("/dictionary").WithHeader("Authorization", bt).WithJSON(dictionary).Expect().Status(http.StatusOK)
}

func TestGetDictionary(t *testing.T) {
	app := web.NewTestApplication(t, new(DictionaryController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	assert.Equal(t, nil, err)
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Get("/dictionary").WithHeader("Authorization", bt).WithQuery("id", 1).Expect().Status(http.StatusOK)
}


func TestDeleteDictionary(t *testing.T) {
	app := web.NewTestApplication(t, new(DictionaryController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	assert.Equal(t, nil, err)
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Delete("/dictionary").WithHeader("Authorization", bt).WithQuery("id", 1).Expect().Status(http.StatusOK)
}
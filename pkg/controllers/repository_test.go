package controllers

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"os"
	"time"
	"fmt"
	"net/http"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hicicd/pkg/scm"
)


func TestRepositoryAddType(t *testing.T) {
	app := web.NewTestApplication(t, new(RepositoryController))
	jwtToken, err := web.GenerateJwtToken(web.JwtMap{
		"url": os.Getenv("SCM_URL"),
		"username": os.Getenv("SCM_USERNAME"),
		"password": os.Getenv("SCM_PASSWORD"), // TODO: token is not working?
		"scmToken": os.Getenv("SCM_TOKEN"),
		"uid": 190,
	}, 1000, time.Hour)
	project := scm.Project{
		ID:  905,
		Ref: "master",
	}
	assert.Equal(t, nil, err)
	bt := fmt.Sprintf("Bearer %v", string(*jwtToken))
	app.Post("/repository/appType").WithHeader("Authorization", bt).WithJSON(project).Expect().Status(http.StatusOK)
}

package controllers

import (
	"testing"
	"os"
	"time"
	"net/http"
	"github.com/iris-contrib/httpexpect"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/auth"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hicicd/pkg/ci"
	"github.com/hidevopsio/hiboot/pkg/utils"
)


func init() {
	utils.ChangeWorkDir("../../")

	userRequest = UserRequest{
		Url:      os.Getenv("SCM_URL"),
		Username: os.Getenv("SCM_USERNAME"),
		Password: os.Getenv("SCM_PASSWORD"),
	}
}


func newTestServer(t *testing.T) *httpexpect.Expect {
	e, err := web.NewTestServer(t, &CicdController{})
	assert.Equal(t, nil, err)

	return e
}

func login(expired int64, unit time.Duration) (*web.Token, error) {
	u := &auth.User{}
	_, _, err := u.Login(userRequest.Url, userRequest.Username, userRequest.Password)
	if err != nil {
		return nil, err
	}
	token, err := web.GenerateJwtToken(web.JwtMap{
		"url":      userRequest.Url,
		"username": userRequest.Username,
		"password": userRequest.Password,
	}, expired, unit)
	return token, err
}

func requestCicdPipeline(e *httpexpect.Expect, token *web.Token, statusCode int, pl *ci.Pipeline) {
	e.Request("POST", "/cicd/run").WithHeader(
		"Authorization", "Bearer "+ string(*token),
	).WithJSON(pl).Expect().Status(statusCode)
}


func TestCicdRunWithExpiredToken(t *testing.T) {
	log.Println("TestCicdRunWithExpiredToken()")

	e := newTestServer(t)

	token, err := login(500, time.Millisecond)
	assert.Equal(t, nil, err)

	if err == nil {
		time.Sleep(1000 * time.Millisecond)

		requestCicdPipeline(e, token, http.StatusUnauthorized, &ci.Pipeline{
			Name:    "java",
			Project: "demo",
			Profile: "dev",
			App:     "hello-world",
		})
	}
}

func TestCicdRunWithoutToken(t *testing.T) {
	log.Println("TestCicdRunWithoutToken()")

	e := newTestServer(t)

	e.Request("POST", "/cicd/run").WithJSON(ci.Pipeline{
		Project: "demo",
		App:     "hello-world",
		Profile: "dev",
		Name:    "java",
	}).Expect().Status(http.StatusUnauthorized)

}

func TestCicdRunJava(t *testing.T) {
	log.Println("TestCicdRun()")

	e := newTestServer(t)

	jwtToken, err := login(24, time.Hour)
	assert.Equal(t, nil, err)

	if err == nil {
		requestCicdPipeline(e, jwtToken, http.StatusOK, &ci.Pipeline{
			Name:    "java",
			Project: "demo",
			Profile: "dev",
			App:     "hello-world",
			//BuildConfigs: ci.BuildConfigs{Skip: true},
			//DeploymentConfigs: ci.DeploymentConfigs{Skip: true},
		})
	}
}

func TestCicdRunNodejs(t *testing.T) {
	log.Println("TestCicdRun()")

	e := newTestServer(t)

	jwtToken, err := login(24, time.Hour)
	assert.Equal(t, nil, err)

	if err == nil {
		requestCicdPipeline(e, jwtToken, http.StatusOK, &ci.Pipeline{
			Name:    "nodejs",
			Project: "demo",
			Profile: "dev",
			App:     "hello-angular",
		})
	}
}

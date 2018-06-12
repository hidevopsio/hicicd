package controllers

import (
	"testing"
	"os"
	"time"
	"net/http"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/auth"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hicicd/pkg/ci"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"fmt"
)


func init() {
	utils.ChangeWorkDir("../../")

	userRequest = UserRequest{
		Url:      os.Getenv("SCM_URL"),
		Username: os.Getenv("SCM_USERNAME"),
		Password: os.Getenv("SCM_PASSWORD"),
	}
}

func login(expired int64, unit time.Duration) (*web.Token, error) {
	u := &auth.User{}
	_, _, _, err := u.Login(userRequest.Url, userRequest.Username, userRequest.Password)
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

func requestCicdPipeline(ta *web.TestApplication, token *web.Token, statusCode int, pl *ci.Pipeline) {
	tk := string(*token)

	log.Println("token: ", tk)

	authToken := fmt.Sprintf("Bearer %v", tk)

	ta.Post("/cicd/run").WithHeader(
		"Authorization", authToken,
	).WithJSON(pl).Expect().Status(statusCode)
}


func TestCicdRunWithExpiredToken(t *testing.T) {
	log.Println("TestCicdRunWithExpiredToken()")

	ta := web.NewTestApplication(t)

	token, err := login(500, time.Millisecond)
	assert.Equal(t, nil, err)

	if err == nil {
		time.Sleep(1000 * time.Millisecond)

		requestCicdPipeline(ta, token, http.StatusUnauthorized, &ci.Pipeline{
			Name:    "java",
			Project: "demo",
			Profile: "dev",
			App:     "hello-world",
		})
	}
}

func TestCicdRunWithoutToken(t *testing.T) {
	log.Println("TestCicdRunWithoutToken()")

	web.NewTestApplication(t).
		Post("/cicd/run").WithJSON(ci.Pipeline{
			Project: "demo",
			App:     "hello-world",
			Profile: "dev",
			Name:    "java",
		}).
		Expect().Status(http.StatusUnauthorized)

}

func TestCicdRunJava(t *testing.T) {
	log.Println("TestCicdRun()")

	ta := web.NewTestApplication(t)

	jwtToken, err := login(24, time.Hour)
	assert.Equal(t, nil, err)

	if err == nil {
		requestCicdPipeline(ta, jwtToken, http.StatusOK, &ci.Pipeline{
			Name:    "java",
			Project: "demo",
			Profile: "test",
			App:     "hello-world",
			Version: "v1",
			BuildConfigs: ci.BuildConfigs{Skip: true},
			//DeploymentConfigs: ci.DeploymentConfigs{Skip: true},
		})
	}
}

func TestCicdRunNodejs(t *testing.T) {
	log.Println("TestCicdRun()")

	ta := web.NewTestApplication(t)

	jwtToken, err := login(24, time.Hour)
	assert.Equal(t, nil, err)

	if err == nil {
		requestCicdPipeline(ta, jwtToken, http.StatusOK, &ci.Pipeline{
			Name:    "nodejs",
			Project: "demo",
			Profile: "dev",
			App:     "hello-angular",
		})
	}
}

var jobChan chan int

func worker(jobChan <- chan int)  {
	for job := range jobChan{
		fmt.Printf("执行任务 %d \n", job)
	}
}

func TestGoChan(t *testing.T)  {
	jobChan = make(chan int, 100)
	//入队
	for i := 1; i <= 10; i++{
		jobChan <- i
	}

	close(jobChan)
	go worker(jobChan)
}



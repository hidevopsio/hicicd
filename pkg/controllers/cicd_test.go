package controllers

import (
	"testing"
	"os"
	"time"
	"net/http"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/auth"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"fmt"
	"github.com/hidevopsio/hicicd/pkg/entity"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
)


func init() {
	io.ChangeWorkDir("../../")

	userRequest = UserRequest{
		Url:      os.Getenv("SCM_URL"),
		Username: os.Getenv("SCM_USERNAME"),
		Password: os.Getenv("SCM_PASSWORD"),
		Uid: os.Getenv("Uid"),
	}
}

func login(expired int64, unit time.Duration) (string, error) {
	u := &auth.User{}
	_, _, _, err := u.Login(userRequest.Url, userRequest.Username, userRequest.Password)
	if err != nil {
		return "", err
	}
	token, err := c.jwtToken.Generate(jwt.Map{
		"url":      userRequest.Url,
		"username": userRequest.Username,
		"password": userRequest.Password, // TODO: token is not working?
		"uid":      userRequest.Uid,
	}, 10, time.Minute)
	return token, err
}

func requestCicdPipeline(ta *web.TestApplication, token string, statusCode int, pl *entity.Pipeline) {
	log.Println("token: ", token)

	//authToken := fmt.Sprintf("Bearer %v", tk)
	//app := web.NewTestApplication(t, new(ConfigMapController))
	//ta.Post("/cicd/run").WithHeader(
	//	"Authorization", authToken,
	//).WithJSON(pl).Expect().Status(statusCode)
}


func TestCicdRunWithExpiredToken(t *testing.T) {
	log.Println("TestCicdRunWithExpiredToken()")

	ta := web.NewTestApplication(t)

	token, err := login(500, time.Millisecond)
	assert.Equal(t, nil, err)

	if err == nil {
		time.Sleep(1000 * time.Millisecond)

		requestCicdPipeline(ta, token, http.StatusUnauthorized, &entity.Pipeline{
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
		Post("/cicd/run").WithJSON(entity.Pipeline{
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
		requestCicdPipeline(ta, jwtToken, http.StatusOK, &entity.Pipeline{
			Name:    "java",
			Project: "demo",
			Profile: "test",
			App:     "hello-world",
			Version: "v1",
			BuildConfigs: entity.BuildConfigs{Enable: true},
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
		requestCicdPipeline(ta, jwtToken, http.StatusOK, &entity.Pipeline{
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



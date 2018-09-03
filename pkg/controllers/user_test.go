package controllers

import (
	"testing"
	"os"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/app/web"
)

var userRequest UserRequest

func init() {

	userRequest = UserRequest{
		Url:      os.Getenv("SCM_URL"),
		Username: os.Getenv("SCM_USERNAME"),
		Password: os.Getenv("SCM_PASSWORD"),
	}
}

func TestUserLogin(t *testing.T) {
	log.Println("TestUserLogin()")

	ta := web.NewTestApplication(t, &UserController{})

	response := ta.Request("POST", "/user/login", ).WithJSON(
		userRequest).Expect().Status(http.StatusOK).JSON().Object()

	response.Value("message").Equal("Success")
}

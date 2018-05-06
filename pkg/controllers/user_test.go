package controllers

import (
	"testing"
	"os"
	"net/http"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
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

	e, err := web.NewTestServer(t, &UserController{})
	assert.Equal(t, nil, err)

	response := e.Request("POST", "/user/login", ).WithJSON(
		userRequest).Expect().Status(http.StatusOK).JSON().Object()
	response.Value("message").Equal("Success")
}

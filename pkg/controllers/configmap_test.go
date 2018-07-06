package controllers

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"net/http"
)

func TestCreate(t *testing.T) {

	ta := web.NewTestApplication(t, &ConfigMapController{})

	response := ta.Request("POST", "/configMap/add", ).WithJSON(
		userRequest).Expect().Status(http.StatusOK).JSON().Object()

	response.Value("message").Equal("Success")


}

package kong

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"os"
	"strings"
)

func TestGet(t *testing.T) {
	api := ApiRequest{
		Name: "Example",
	}
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	a, err := api.Get(baseUrl)
	assert.Equal(t, nil, err)
	log.Info("result: ", a)
}

func TestPost(t *testing.T) {
	upstreamUrl := "http://hello-world-demo-dev.apps.oc.com"
	namespace := "demo-dev"
	app := "hello-world"
	uris := "/"+namespace + "-" + app
	uris = strings.Replace(uris, "-", "/", -1)
	//host := os.Getenv("KONG_HOST")
	apiRequest := &ApiRequest{
		Name:                   app,
		Hosts:                  []string{"kong-proxy-kong-gateway.apps.oc.com"},
		Uris:                   []string{uris},
		UpstreamURL:            upstreamUrl,
		StripUri:               true,
		PreserveHost:           false,
		Retries:                5,
		UpstreamConnectTimeout: 6000,
		UpstreamSendTimeout:    6000,
		UpstreamReadTimeout:    6000,
		HttpsOnly:              false,
		HttpIfTerminated:       true,
	}
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	 err := apiRequest.Post(baseUrl)
	assert.Equal(t, nil, err)
}

func TestDelete(t *testing.T) {
	api := ApiRequest{
		Name: "Example",
	}
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	err := api.Delete(baseUrl)
	assert.Equal(t, nil, err)
}
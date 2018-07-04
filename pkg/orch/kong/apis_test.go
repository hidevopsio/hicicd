package kong

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"strings"
	"os"
)

func TestGet(t *testing.T) {
	api := ApiRequest{
		Name: "hello-world",
	}
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	baseUrl = strings.Replace(baseUrl, "${profile}", "stage", -1)
	a, err := api.Get(baseUrl)
	assert.Equal(t, nil, err)
	log.Info("result: ", a)
}

func TestPost(t *testing.T) {
	upstreamUrl := "hello-world.demo-stage:8080"
	app := "hello-world"
	uris := "/demo/hello/world"
	log.Debug("Pipeline.CreateKongGateway()")
	host := os.Getenv("KONG_HOST")
	host = strings.Replace(host, "${profile}", "stage", -1)
	hosts := strings.Split(host, ",")
	apiRequest := &ApiRequest{
		Name:                   app,
		Hosts:                  hosts,
		Uris:                   []string{uris},
		UpstreamURL:            "http://" + upstreamUrl,
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
	baseUrl = strings.Replace(baseUrl, "${profile}", "stage", -1)
	err := apiRequest.Post(baseUrl)
	assert.Equal(t, nil, err)

}

func TestDelete(t *testing.T) {
	api := ApiRequest{
		Name: "hello-world",
	}
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	baseUrl = strings.Replace(baseUrl, "${profile}", "stage", -1)
	err := api.Delete(baseUrl)
	assert.Equal(t, nil, err)
}
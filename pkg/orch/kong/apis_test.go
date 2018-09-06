package kong

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"strings"
	"os"
	"github.com/kevholditch/gokong"
)

func TestApisGet(t *testing.T) {
	name := "hello-world"
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	baseUrl = strings.Replace(baseUrl, "${profile}", "stage", -1)
	config := &gokong.Config{
		HostAddress: baseUrl,
	}
	reult, err := gokong.NewClient(config).Apis().GetByName(name)
	assert.Equal(t, nil, err)
	log.Info("result: ", reult)
}


func TestApisDelete(t *testing.T) {
	name := "hello-world"
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	baseUrl = strings.Replace(baseUrl, "${profile}", "stage", -1)
	config := &gokong.Config{
		HostAddress: baseUrl,
	}
	err := gokong.NewClient(config).Apis().DeleteByName(name)
	assert.Equal(t, nil, err)
}

func Test_ApisCreate(t *testing.T) {
	upstreamUrl := "hello-world.demo-stage:8080"
	app := "hello-world"
	uris := "/demo/hello/world"
	log.Debug("Pipeline.CreateKongGateway()")
	host := os.Getenv("KONG_HOST")
	host = strings.Replace(host, "${profile}", "stage", -1)
	hosts := strings.Split(host, ",")
	apiRequest := &gokong.ApiRequest{
		Name:                   app,
		Hosts:                  hosts,
		Uris:                   []string{uris},
		UpstreamUrl:            "http://" + upstreamUrl,
		StripUri:               false,
		PreserveHost:           true,
		Retries:                "3",
		UpstreamConnectTimeout: 1000,
		UpstreamSendTimeout:    2000,
		UpstreamReadTimeout:    3000,
		HttpsOnly:              false,
		HttpIfTerminated:       true,
	}
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	baseUrl = strings.Replace(baseUrl, "${profile}", "stage", -1)
	config := &gokong.Config{
		HostAddress: baseUrl,
	}
	result, err := gokong.NewClient(config).Apis().Create(apiRequest)
	assert.Equal(t, nil, err)
	assert.Equal(t, apiRequest.Name, result.Name)
	assert.Equal(t, apiRequest.Hosts, result.Hosts)
	assert.Equal(t, apiRequest.Uris, result.Uris)
	assert.Equal(t, apiRequest.Methods, result.Methods)
	assert.Equal(t, apiRequest.UpstreamUrl, result.UpstreamUrl)
	assert.Equal(t, apiRequest.StripUri, result.StripUri)
	assert.Equal(t, apiRequest.PreserveHost, result.PreserveHost)
	assert.Equal(t, apiRequest.UpstreamConnectTimeout, result.UpstreamConnectTimeout)
	assert.Equal(t, apiRequest.UpstreamSendTimeout, result.UpstreamSendTimeout)
	assert.Equal(t, apiRequest.UpstreamReadTimeout, result.UpstreamReadTimeout)
	assert.Equal(t, apiRequest.HttpsOnly, result.HttpsOnly)
	assert.Equal(t, apiRequest.HttpIfTerminated, result.HttpIfTerminated)
}
package kong

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"strings"
)

func TestGet(t *testing.T) {
	api := ApiRequest{
		Name: "moses-comment-consumer",
	}
	baseUrl := "http://kong-admin-kong-gateway-stage.apps.cloud.vpclub.cn"
	a, err := api.Get(baseUrl)
	assert.Equal(t, nil, err)
	log.Info("result: ", a)
}

func TestPost(t *testing.T) {
	upstreamUrl := "comment-consumer-moses-stage.apps.cloud.vpclub.cn"
	app := "moses-comment-consumer"
	uris := "/moses/comment/consumer"
	uris = strings.Replace(uris, "-", "/", -1)
	log.Debug("Pipeline.CreateKongGateway()")
	uris = strings.Replace(uris, "-", "/", -1)
	host := "stagecould.vpclub.cn"
	apiRequest := &ApiRequest{
		Name:                   app,
		Hosts:                  []string{host},
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
	baseUrl := "http://kong-admin-kong-gateway-stage.apps.cloud.vpclub.cn"
	baseUrl = strings.Replace(baseUrl, "${profile}", "stage", -1)
	err := apiRequest.Post(baseUrl)
	assert.Equal(t, nil, err)

}

func TestDelete(t *testing.T) {
	api := ApiRequest{
		Name: "comment-consumer-moses-stage-legacy",
	}
	baseUrl := "http://kong-admin-kong-gateway-stage.apps.cloud.vpclub.cn"
	err := api.Delete(baseUrl)
	assert.Equal(t, nil, err)
}
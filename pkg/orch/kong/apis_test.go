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
	upstreamUrl := "hiweb-hidevopsio.apps.cloud.vpclub.cn"
	app := "hiweb"
	uris := "/hidevopsio/hiweb"
	uris = strings.Replace(uris, "-", "/", -1)
	log.Debug("Pipeline.CreateKongGateway()")
	uris = strings.Replace(uris, "-", "/", -1)
	host := "hicloud.vpclub.cn"
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
		Name: "express-consumer",
	}
	baseUrl := "http://kong-admin-kong-gateway-prod.apps.cloud.vpclub.cn"
	err := api.Delete(baseUrl)
	assert.Equal(t, nil, err)
}
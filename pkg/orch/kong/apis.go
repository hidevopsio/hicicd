package kong

import (
	"github.com/nccurry/go-kong/kong"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/jinzhu/copier"
)

type ApiRequest struct {
	UpstreamURL            string   `json:"upstream_url,omitempty"`
	RequestPath            string   `json:"request_path,omitempty"`
	ID                     string   `json:"id,omitempty"`
	CreatedAt              int64    `json:"created_at,omitempty"`
	PreserveHost           bool     `json:"preserve_host,omitempty"`
	Name                   string   `json:"name,omitempty"`
	Hosts                  []string `json:"hosts,omitempty"`
	Uris                   []string `json:"uris"`
	StripUri               bool     `json:"strip_uri"`
	Retries                int      `json:"retries"`
	UpstreamConnectTimeout int      `json:"upstream_connect_timeout"`
	UpstreamSendTimeout    int      `json:"upstream_send_timeout"`
	UpstreamReadTimeout    int      `json:"upstream_read_timeout"`
	HttpsOnly              bool      `json:"https_only"`
	HttpIfTerminated       bool      `json:"http_if_terminated"`
}

func (a *ApiRequest) Get(baseUrl string) (*kong.Api, error) {
	log.Debug("kong get api :")
	c, err := kong.NewClient(&http.Client{}, baseUrl)
	if err != nil {
		log.Error("api get err :", err)
		return nil, err
	}
	api, _, err := c.Apis.Get(a.Name)
	log.Info("kong get api response:", err)
	return api, err
}

func (a *ApiRequest) Post(baseUrl string) error {
	log.Debug("kong post()",)
	c, err := kong.NewClient(&http.Client{}, baseUrl)
	if err != nil {
		log.Error("api post err :", err)
		return err
	}
	api, err := a.Get(baseUrl)
	if err == nil {
		if api.Name == a.Name {
			return nil
		}
		 err := a.Update(baseUrl)
		return err
	}
	apiRequest := &kong.ApiRequest{}
	copier.Copy(apiRequest, a)
	r, err := c.Apis.Post(apiRequest)
	log.Info("kong post api call back:", r)
	return err
}

func (a *ApiRequest) Delete(baseUrl string) error  {
	log.Debug("kong get api :", a)
	c, err := kong.NewClient(&http.Client{}, baseUrl)
	if err != nil {
		log.Error("api get err :", err)
		return err
	}
	_, err = c.Apis.Delete(a.Name)
	log.Info("kong get api response:", err)
	return err
}

func (a *ApiRequest) Update(baseUrl string) error {
	log.Debug("kong get api :", a)
	c, err := kong.NewClient(&http.Client{}, baseUrl)
	if err != nil {
		log.Error("api get err :", err)
		return err
	}
	kongApi := &kong.ApiRequest{}
	copier.Copy(kongApi, a)
	_, err = c.Apis.Patch(kongApi)
	log.Info("kong get api response:", err)
	return err

}
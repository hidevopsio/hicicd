package auth

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"crypto/tls"
	"net/http"
	"io/ioutil"
	"github.com/hidevopsio/hicicd/pkg/utils"
)

func TestGetAccessToken(t *testing.T)  {
	code := "e480550bd1f668633581e3a5cf180b7fcf97eb4faec0f6ee002634dcac9663a3"
	//s := url.QueryEscape(CallbackUrl)
	session := NewClient(BaseUrl, AccessTokenUrl, ApplicationId, CallbackUrl, Secret)
	_, err := session.GetAccessToken(code)
	assert.Equal(t, err, nil)
}

func TestHttpsConnection(t *testing.T) {
	transport := &utils.Transport{
		ConnectTimeout: 1 * time.Second,
		RequestTimeout: 2 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defer transport.Close()
	client := &http.Client{Transport: transport}

	req, _ := http.NewRequest("GET", "https://httpbin.org/ip", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("1st request failed - %s", err.Error())
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("1st failed to read body - %s", err.Error())
	}
	resp.Body.Close()

	req2, _ := http.NewRequest("GET", "https://httpbin.org/delay/5", nil)
	_, err = client.Do(req2)
	if err == nil {
		t.Fatalf("HTTPS request should have timed out")
	}
}
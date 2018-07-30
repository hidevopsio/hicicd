package fake

import (
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	client *http.Client

	baseURL *url.URL

	tokenType tokenType

	token string

	UserAgent string

	session *FakeSession
}
type tokenType int
const (
	privateToken tokenType = iota
	oAuthToken
)

func NewClient(httpClient *http.Client, token string) *Client  {
	return newClient(httpClient, privateToken, token)
}

func newClient(httpClient *http.Client, tokenType tokenType, token string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{client: httpClient, tokenType: tokenType, token: token, UserAgent: ""}
	c.session = &FakeSession{client: c}
	return c
}

func (c *Client) SetBaseURL(urlStr string) error {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.baseURL, err = url.Parse(urlStr)
	return err
}
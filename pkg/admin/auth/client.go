package auth

import (
	"time"
	"crypto/tls"
	"net/http"
	"github.com/hidevopsio/hicicd/pkg/utils"
	"io"
	"encoding/json"
)

func Client(method, baseUrl string, v interface{})  (*http.Response, error){
	transport := &utils.Transport{
		ConnectTimeout: 1 * time.Second,
		RequestTimeout: 2 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defer transport.Close()
	client := &http.Client{Transport: transport}

	req2, _ := http.NewRequest(method, baseUrl, nil)
	resp, err := client.Do(req2)
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	resp.Body.Close()
	return resp, err
}
// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitlab

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/stretchr/testify/assert"
	"os"
	"net/http"
	"net/http/httptest"
	"github.com/xanzy/go-gitlab"
)

func init()  {
	log.SetLevel(log.DebugLevel)
}

func TestUserGet(t *testing.T)  {
	mux, test, client := setup()
	baseUrl :=  os.Getenv("SCM_URL")
	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	assert.Equal(t, username, gs.Username)
	log.Info(gs)
}

func setup() (*http.ServeMux, *httptest.Server, *gitlab.Client) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// client is the Gitlab client being tested.
	client := gitlab.NewClient(nil, "")
	client.SetBaseURL(server.URL)

	return mux, server, client
}
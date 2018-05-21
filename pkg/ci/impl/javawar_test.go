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

package impl

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/ci"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hicicd/pkg/scm/gitlab"
)

func init()  {
	log.SetLevel(log.DebugLevel)
}

func TestJavaWarPipeline(t *testing.T)  {

	log.Debug("Test Java War Pipeline")

	java := &JavaWarPipeline{}

	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")

	java.Init(&ci.Pipeline{
		Name: "java-war",
		Profile: "dev",
		App: "hello-war",
		Project: "demo",
		Scm: ci.Scm{Url: os.Getenv("SCM_URL")},
	})
	baseUrl :=  os.Getenv("SCM_URL")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	log.Debug(gs)
	assert.Equal(t, username, gs.Username)
	err = java.Run(username, password,gs.PrivateToken, false)
	assert.Equal(t, nil, err)
}

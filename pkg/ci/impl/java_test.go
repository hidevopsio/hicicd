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
	"reflect"
	"fmt"
	"github.com/hidevopsio/hicicd/pkg/scm/gitlab"
)

func init()  {
	log.SetLevel(log.DebugLevel)
}

func TestJavaPipeline(t *testing.T)  {

	log.Debug("Test Java Pipeline")

	java := &JavaPipeline{}

	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	baseUrl :=  os.Getenv("SCM_URL")
	java.Init(&ci.Pipeline{
		Name: "java",
		Profile: "test",
		App: "hello-world",
		Project: "demo",
		Version: "v2",
		Scm: ci.Scm{
			Url: baseUrl,
			Ref: "v2",
		},
		BuildConfigs: ci.BuildConfigs{
			Skip: false,
		},
		DeploymentConfigs: ci.DeploymentConfigs{
			ForceUpdate: true,
		},
	})
	gs := new(gitlab.Session)
	err := gs.GetSession(baseUrl, username, password)
	assert.Equal(t, nil, err)
	log.Debug(gs)
	assert.Equal(t, username, gs.Username)
	err = java.Run(username, password, gs.PrivateToken, gs.ID, false)
	assert.Equal(t, nil, err)
}


type Book struct {
	Id    int
	Title string
	Price float32
	Authors []string
}

func TestIterateStruct(t *testing.T) {
	book := Book{Id: 12, Title: "test"}
	e := reflect.ValueOf(&book).Elem()

	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varType := e.Type().Field(i).Type
		varValue := e.Field(i).Interface()
		fmt.Printf("%v %v %v\n", varName,varType,varValue)
	}
}

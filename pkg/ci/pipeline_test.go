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

package ci

import (
	"testing"
	"github.com/magiconair/properties/assert"
)

func TestPipelineInit(t *testing.T) {

}

func TestPipelineCreateProject(t *testing.T) {
	name := "demo-test"
	pl := &Pipeline{
		Name: name,
	}
	err := pl.CreateProject()
	assert.Equal(t, nil, err)
}

func TestPipelineCreateRoleBinding(t *testing.T) {
	username := "chulei"
	p := &Pipeline{
		Name: "test",
		Namespace: "demo-dev",
	}
	metaName := "admin"
	roleRefName := "master"
	err := p.CreateRoleBinding(username, metaName, roleRefName)
	assert.Equal(t, nil, err)
}
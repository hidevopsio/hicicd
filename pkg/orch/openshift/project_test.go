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

package openshift

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestProjectLit(t *testing.T) {
	projectName := "project-crud"
	project, err := NewProject(projectName, projectName, "project for testing")
	assert.Equal(t, nil, err)

	pl, err := project.List()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(pl.Items))
	log.Debugf("There are %d projects in the cluster", len(pl.Items))

	for i, p := range pl.Items {
		log.Debugf("index %d: project: %s", i, p.Name)
	}
}

func TestProjectCrud(t *testing.T) {
	projectName := "project-crud"
	project, err := NewProject(projectName, projectName, "project for testing")
	assert.Equal(t, nil, err)

	// create project
	p, err := project.Create()
	assert.Equal(t, nil, err)
	assert.Equal(t, projectName, p.Name)

	// read project
	p, err = project.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, projectName, p.Name)

	// delete project
	err = project.Delete()
	assert.Equal(t, nil, err)

}


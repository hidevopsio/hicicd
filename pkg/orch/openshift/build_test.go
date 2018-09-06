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
	"os"
	"github.com/hidevopsio/hiboot/pkg/system"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestBuildCreation(t *testing.T) {
	log.Debug("TestBuildCreate()")

	// put below configs in yaml file
	project := "demo"
	profile := "stage"
	namespace := project + "-" + profile
	appName := "hello-world"
	scmUrl := os.Getenv("SCM_URL") + "/" + project + "/" + appName + ".git"
	scmRef := "master"
	secret := "test-secret"
	version := "v1"
	s2iImageStream := "s2i-java:latest"
	repoUrl := os.Getenv("MAVEN_MIRROR_URL")
	env := []system.Env{
		{
			Name: "BUILD_SCRIPT",
			Value: "mvn clean package -Dmaven.test.skip=true -Djava.net.preferIPv4Stack=true",
		},
		{
			Name: "MAVEN_MIRROR_URL",
			Value: os.Getenv("MAVEN_MIRROR_URL"),
		},
	}
	log.Debug(repoUrl)
	log.Debug(scmUrl)

	log.Debugf("workDir: %v", os.Getenv("PWD"))

	buildConfig, err := NewBuildConfig(namespace, appName, scmUrl, scmRef, secret, version, s2iImageStream, false)
	assert.Equal(t, nil, err)

	bc, err := buildConfig.Create()
	assert.Equal(t, nil, err)
	assert.Equal(t, appName, bc.Name)

	// Get build config
	bc, err = buildConfig.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, appName, bc.Name)

	// Build image stream
	build, err := buildConfig.Build(env)
	assert.Equal(t, nil, err)
	assert.Contains(t, build.Name, appName)

	log.Debug("End of build test")
}



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

type PipelineInterface interface {
	Init(pl *Pipeline)
	CreateSecret(username, password string, isToken bool) (string, error)
	Build(secret string, completedHandler func() error) error
	RunUnitTest() error
	RunIntegrationTest() error
	Analysis() error
	CreateDeploymentConfig(force bool, injectFn func(in interface{}) (interface{}, error)) error
	Deploy() error
	CreateService() error
	CreateRoute() (string, error)
	Run(username, password, scmToken string, uid int, isToken bool) error
	InitProject() error
	CreateRoleBinding(username, metaName, roleRefName string) error
	CreateProject() error
	InstantiateDeploymentConfig() error
	CreateImageStreamTag() error
}

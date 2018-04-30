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

package istio

import (
	"istio.io/istio/pilot/pkg/kube/inject"
	"istio.io/istio/pilot/pkg/model"
	"k8s.io/apimachinery/pkg/runtime"
)

type Injector struct {
	Hub string
	Tag string
	Version string
	DebugMode bool
}

func (i *Injector) Inject(in runtime.Object) (interface{}, error)  {
	mesh := model.DefaultMeshConfig()
	params := &inject.Params{
		InitImage:           inject.InitImageName(i.Hub, i.Tag, i.DebugMode),
		ProxyImage:          inject.ProxyImageName(i.Hub, i.Tag, i.DebugMode),
		ImagePullPolicy:     "IfNotPresent",
		Verbosity:           inject.DefaultVerbosity,
		SidecarProxyUID:     inject.DefaultSidecarProxyUID,
		Version:             i.Version,
		EnableCoreDump:      false,
		Mesh:                &mesh,
		DebugMode:           i.DebugMode,
	}

	sidecarTemplate, err := inject.GenerateTemplateFromParams(params)
	if err != nil {
		return nil, err
	}

	out, err := inject.IntoObject(sidecarTemplate, &mesh, in)
	return out, err
}
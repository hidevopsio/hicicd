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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/openshift/api/apps/v1"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"istio.io/istio/pilot/pkg/kube/inject"
	"istio.io/istio/pilot/pkg/model"
)

const (
	unitTestHub = "docker.io/istio"
	unitTestTag = "unittest"
	name = "foo"
)

func getDeploymentConfig() *v1.DeploymentConfig  {
	return &v1.DeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: v1.DeploymentConfigSpec{
			Replicas: 1,

			Selector: map[string]string{
				"app": name,
			},

			Strategy: v1.DeploymentStrategy{
				Type: v1.DeploymentStrategyTypeRolling,
			},

			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
					Labels: map[string]string{
						"app":  name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							//Env:             e,
							Image:           " ",
							ImagePullPolicy: corev1.PullAlways,
							Name:            name,
							//Ports:           p,
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command : []string{
											"curl",
											"--silent",
											"--show-error",
											"--fail",
											"http://localhost:8080/health",
										},
									},
								},
								InitialDelaySeconds: 10,
								TimeoutSeconds:      1,
								PeriodSeconds:       5,
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command : []string{
											"curl",
											"--silent",
											"--show-error",
											"--fail",
											"http://localhost:8080/health",
										},
									},
								},
								InitialDelaySeconds: 20,
								TimeoutSeconds:      1,
								PeriodSeconds:       5,
							},
						},
					},
					DNSPolicy:     corev1.DNSClusterFirst,
					RestartPolicy: corev1.RestartPolicyAlways,
					SchedulerName: "default-scheduler",
				},
			},
			Test: false,
			Triggers: v1.DeploymentTriggerPolicies{
				{
					Type: v1.DeploymentTriggerOnImageChange,
					ImageChangeParams: &v1.DeploymentTriggerImageChangeParams{
						Automatic: true,
						ContainerNames: []string{
							name,
						},
						From: corev1.ObjectReference{
							Kind:      "ImageStreamTag",
							Name:      name + ":" + "latest",
							Namespace: "demo-dev",
						},
					},
				},
			},
		},
	}
}



func TestIntoObject(t *testing.T) {

	cfg := getDeploymentConfig()
	log.Print(cfg)
	assert.Equal(t, 1, len(cfg.Spec.Template.Spec.Containers))

	mesh := model.DefaultMeshConfig()

	params := &inject.Params{
		InitImage:           inject.InitImageName(unitTestHub, unitTestTag, true),
		ProxyImage:          inject.ProxyImageName(unitTestHub, unitTestTag, true),
		ImagePullPolicy:     "IfNotPresent",
		Verbosity:           inject.DefaultVerbosity,
		SidecarProxyUID:     inject.DefaultSidecarProxyUID,
		Version:             "v1",
		EnableCoreDump:      false,
		Mesh:                &mesh,
		DebugMode:           true,
		IncludeIPRanges:     inject.DefaultIncludeIPRanges,
		ExcludeIPRanges:     "",
		IncludeInboundPorts: inject.DefaultIncludeInboundPorts,
		ExcludeInboundPorts: "",
	}

	sidecarTemplate, err := inject.GenerateTemplateFromParams(params)

	out, err := inject.IntoObject(sidecarTemplate, &mesh, cfg)
	assert.Equal(t, nil, err)

	dc := out.(*v1.DeploymentConfig)
	assert.Equal(t, 2, len(dc.Spec.Template.Spec.Containers))

	log.Print(dc)
}


func TestInjectorInject(t *testing.T) {

	injector := &Injector{
		Version: "0.7.1",
		Namespace: "istio-system",
		MeshConfigMapName: "istio",
		InjectConfigMapName: "istio-inject",
		DebugMode: false,
		SidecarProxyUID: uint64(1337),
		Verbosity: 2,
		ImagePullPolicy: "IfNotPresent",
		IncludeIPRanges: "*",
		IncludeInboundPorts: "*",
	}

	cfg := getDeploymentConfig()
	out, err := injector.Inject(cfg)
	assert.Equal(t, nil, err)

	dc := out.(*v1.DeploymentConfig)
	assert.Equal(t, 2, len(dc.Spec.Template.Spec.Containers))

	log.Print(dc)
}

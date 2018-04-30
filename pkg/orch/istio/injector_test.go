// Copyright 2018 Istio Authors
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
	"istio.io/istio/pilot/pkg/model"
	"github.com/stretchr/testify/assert"
	"istio.io/istio/pilot/pkg/kube/inject"

	"github.com/hidevopsio/hiboot/pkg/log"
)


const (
	// This is the hub to expect in
	// platform/kube/inject/testdata/frontend.yaml.injected and the
	// other .injected "want" YAMLs
	unitTestHub = "docker.io/istio"

	// Tag name should be kept in sync with value in
	// platform/kube/inject/refresh.sh
	unitTestTag = "unittest"
)

func TestIntoObject(t *testing.T) {
	debugMode := true
	mesh := model.DefaultMeshConfig()
	params := &inject.Params{
		InitImage:           inject.InitImageName(unitTestHub, unitTestTag, debugMode),
		ProxyImage:          inject.ProxyImageName(unitTestHub, unitTestTag, debugMode),
		ImagePullPolicy:     "IfNotPresent",
		Verbosity:           inject.DefaultVerbosity,
		SidecarProxyUID:     inject.DefaultSidecarProxyUID,
		Version:             "12345678",
		EnableCoreDump:      false,
		Mesh:                &mesh,
		DebugMode:           debugMode,
	}

	sidecarTemplate, err := inject.GenerateTemplateFromParams(params)
	assert.Equal(t, nil, err)
	log.Debug(sidecarTemplate)
	if err != nil {
		t.Fatalf("GenerateTemplateFromParams(%v) failed: %v", params, err)
	}
	name := "foo"
	cfg := &v1.DeploymentConfig{
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
	log.Print(cfg)
	//
	//out, err := inject.IntoObject(sidecarTemplate, &mesh, cfg)
	//assert.Equal(t, nil, err)
	//
	//dc := out.(*v1.DeploymentConfig)
	//
	//log.Print(dc)
}
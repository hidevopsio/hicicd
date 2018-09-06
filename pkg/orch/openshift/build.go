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
	"github.com/openshift/api/build/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	"github.com/hidevopsio/hiboot/pkg/log"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/watch"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/jinzhu/copier"
	"fmt"
	"github.com/hidevopsio/hicicd/pkg/orch"
	imagev1 "github.com/openshift/api/image/v1"
)

type Scm struct {
	Url    string
	Ref    string
	Secret string
}

type BuildConfig struct {
	Name      string
	Namespace string
	Scm       Scm
	Version   string
	// use NewFrom when creating new buildConfig
	From corev1.ObjectReference

	BuildConfigs buildv1.BuildConfigInterface
	Builds       buildv1.BuildInterface
}

// @Title NewBuildConfig
// @Description Create new BuildConfig Instance
// @Param namespace, appName, gitUrl, imageTag, s2iImageStream string
// @Return *BuildConfig, error
func NewBuildConfig(namespace, name, scmUrl, scmRef, scmSecret, version, s2iImageStream string, rebuild bool) (*BuildConfig, error) {

	log.Debug("NewBuildConfig()")

	// TODO: for the sake of decoupling, the image stream creation should be here or not?
	// create imagestream
	var err error
	imageStream, err := NewImageStream(name, namespace)
	if err != nil {
		return nil, err
	}

	var from corev1.ObjectReference
	var is *imagev1.ImageStream
	if !rebuild {
		is, err = imageStream.Get()
		// the images stream is exist with 0 tags, then delete it
		if len(is.Status.Tags) == 0 {
			imageStream.Delete()
			is, err = imageStream.Get()
		}
	}

	// create new images stream if it is not found
	if errors.IsNotFound(err) || rebuild{
		_, err := imageStream.Create(version)
		if err != nil {
			return nil, err
		}
		from = corev1.ObjectReference{
			Kind:      "ImageStreamTag",
			Name:      s2iImageStream,
			Namespace: "openshift",
		}
	} else {
		from = corev1.ObjectReference{
			Kind:      "ImageStreamTag",
			Name:      name + ":" + is.Status.Tags[0].Tag,
			Namespace: namespace,
		}
	}

	clientSet, err := buildv1.NewForConfig(orch.Config)
	buildConfig := &BuildConfig{
		BuildConfigs: clientSet.BuildConfigs(namespace),
		Builds:       clientSet.Builds(namespace),

		From:      from,
		Name:      name,
		Namespace: namespace,
		Scm: Scm{
			Url:    scmUrl,
			Ref:    scmRef,
			Secret: scmSecret,
		},

		Version: version,
	}
	return buildConfig, err
}

// @Title Create
// @Description Create new BuildConfig
// @Param
// @Return *v1.BuildConfig, error
func (b *BuildConfig) Create() (*v1.BuildConfig, error) {
	log.Debug("BuildConfig.Create()")

	// buildConfig
	buildConfig := &v1.BuildConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: b.Name,
			Labels: map[string]string{
				"app": b.Name,
			},
		},
		Spec: v1.BuildConfigSpec{
			// The runPolicy field controls whether builds created from this build configuration can be run simultaneously.
			// The default value is Serial, which means new builds will run sequentially, not simultaneously.
			RunPolicy: v1.BuildRunPolicy("Serial"),
			CommonSpec: v1.CommonSpec{

				Source: v1.BuildSource{
					Type: v1.BuildSourceType(v1.BuildSourceGit),
					Git: &v1.GitBuildSource{
						URI: b.Scm.Url,
						Ref: b.Scm.Ref,
					},
					SourceSecret: &corev1.LocalObjectReference{
						Name: b.Scm.Secret,
					},
				},
				Strategy: v1.BuildStrategy{
					Type: v1.BuildStrategyType(v1.SourceBuildStrategyType),
					SourceStrategy: &v1.SourceBuildStrategy{
						From: b.From,
					},
				},
				Output: v1.BuildOutput{
					To: &corev1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: b.Name + ":" + b.Version,
					},
				},
			},
		},
	}

	bc, err := b.BuildConfigs.Get(b.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		bc, err = b.BuildConfigs.Create(buildConfig)
		if nil == err {
			log.Infof("Created BuildConfig %v", bc.Name)
		}
	} else {
		buildConfig.ResourceVersion = bc.ResourceVersion
		bc, err = b.BuildConfigs.Update(buildConfig)
		if nil == err {
			log.Infof("Updated BuildConfig %v", bc.Name)
		}
	}

	return bc, err
}

// @Title Get
// @Description Get BuildConfig
// @Param
// @Return *v1.BuildConfig, error
func (b *BuildConfig) Get() (*v1.BuildConfig, error) {
	log.Debug("BuildConfig.Get()")
	return b.BuildConfigs.Get(b.Name, metav1.GetOptions{})
}

// @Title Delete
// @Description Delete BuildConfig
// @Param
// @Return error
func (b *BuildConfig) Delete() error {
	log.Debug("BuildConfig.Delet()")
	return b.BuildConfigs.Delete(b.Name, &metav1.DeleteOptions{})
}

// @Title Build
// @Description Start build according to previous build config settings, it will produce new image build
// @Param repo string, buildCmd string
// @Return *v1.Build, error
func (b *BuildConfig) Build(env []system.Env) (*v1.Build, error) {
	log.Debug("BuildConfig.Build()")

	e := make([]corev1.EnvVar, 0)
	copier.Copy(&e, env)

	incremental := false
	buildTriggerCauseManualMsg := "Manually triggered"
	buildRequest := v1.BuildRequest{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BuildRequest",
			APIVersion: "build.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: b.Name,
			Labels: map[string]string{
				"app":     b.Name,
				"version": b.Version,
			},
		},
		TriggeredBy: append([]v1.BuildTriggerCause{},
			v1.BuildTriggerCause{
				Message: buildTriggerCauseManualMsg,
			},
		),
		SourceStrategyOptions: &v1.SourceStrategyOptions{
			Incremental: &incremental,
		},
		Env:  e,
		From: &b.From,
	}

	build, err := b.BuildConfigs.Instantiate(b.Name, &buildRequest)
	if nil != err {
		log.Error("b.BuildConfigs.Instantiate err", err)
		return nil, err
	}
	log.Infof("Instantiated Build %v", build.Name)
	return build, err
}

func (b *BuildConfig) Watch(build *v1.Build, completedHandler func() error) error {
	w, err := b.Builds.Watch(metav1.ListOptions{
		LabelSelector: "app=" + b.Name,
		Watch:         true,
	})

	if nil != err {
		log.Error("BuildConfig.Watch err ", err)
		return err
	}

	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				log.Info("resultChan: ", ok)
				return fmt.Errorf("resultChan: %v", ok)
			}
			switch event.Type {
			case watch.Added:
				bld := event.Object.(*v1.Build)
				log.Info("Added new build ", bld.Name)
			case watch.Modified:

				bld := event.Object.(*v1.Build)
				if bld.Name == build.Name {
					//log.Info("Modified: ", event.Object)
					log.Debugf("bld.Status.Phase: %v", bld.Status.Phase)
					switch bld.Status.Phase {
					case v1.BuildPhaseComplete:
						log.Info("bld.Status.Phase", bld.Status.Phase)
						var err error
						if nil != completedHandler {
							err = completedHandler()
						}
						w.Stop()
						log.Error("bld.Status.Phase completedHandler", err)
						return err
					case v1.BuildPhaseError, v1.BuildPhaseCancelled, v1.BuildPhaseFailed:
						w.Stop()
						log.Error("bld.Status.Phase BuildPhaseError", fmt.Errorf(bld.Status.Message))
						return fmt.Errorf(bld.Status.Message)

					}
				}

			case watch.Deleted:
				log.Info("Deleted: ", event.Object)
			default:
				log.Error("Failed")
			}
		}
	}
	log.Info("build.watch", err)
	return err
}

// @Title GetBuild
// @Description Get current build
// @Param
// @Return *v1.Build, error
func (b *BuildConfig) GetBuild() (*v1.Build, error) {
	log.Debug("BuildConfig.GetBuild()")
	return b.Builds.Get(b.Name, metav1.GetOptions{})
}

// @Title GetBuildStatus
// @Description Get current build status
// @Param
// @Return v1.BuildPhase, error
func (b *BuildConfig) GetBuildStatus() (v1.BuildPhase, error) {
	log.Debug("BuildConfig.GetBuildStatus()")
	build, err := b.GetBuild()
	return build.Status.Phase, err
}

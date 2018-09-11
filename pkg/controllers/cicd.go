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

package controllers

import (
			"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"os"
	"strings"
	"github.com/hidevopsio/hicicd/pkg/service"
	"github.com/hidevopsio/hicicd/pkg/entity"
	"strconv"
)

type PipelineResponse struct {
	model.Response
}

// Operations about object
type PipelineController struct {
	BaseController
	remoteService   *service.RemoteDeploymentConfigsService
}

func (p *PipelineController) Init(remoteService *service.RemoteDeploymentConfigsService) {
	p.remoteService = remoteService
}

const (
	ScmUrl      = "url"
	ScmUsername = "username"
	ScmPassword = "password"
	ScmToken    = "scmToken"
	ScmUid      = "uid"
)

func init() {
	web.RestController(new(PipelineController))
}

func (p *PipelineController) Before(ctx *web.Context) {
	p.BaseController.Before(ctx)
}

// @Title Deploy
// @Description deploy application by the pipeline
// @Param	body
// @Success 200 {string}
// @Failure 403 body is empty
// @router / [post]
func (p *PipelineController) PostRun(ctx *web.Context) {
	log.Debug("Cicd Controller.Run()")
	var pl entity.Pipeline
	err := ctx.RequestBody(&pl)
	if err != nil {
		return
	}
	message := "success"
	if err == nil {
		pipelineService := &service.PipelineService{}
		pipelineService.Initialize(&pl, p.JwtProperty("url"))
		if pipelineService.DeploymentConfigs.RemoteEnable {
			err = p.RemoteDeploy(pipelineService)
		} else {
			go func() {
				uid, _ := strconv.Atoi(p.JwtProperty("uid"))
				err = pipelineService.Run(p.JwtProperty("username"), p.JwtProperty("password"), p.JwtProperty("scmToken"), uid, false)
				if err != nil {
					message = err.Error()
				}
			}()
		}
	} else {
		message = "failed, " + err.Error()
	}
	bc := os.Getenv("BUILD_CONSOLE")
	bc = strings.Replace(bc, "${project}", pl.Project, -1)
	bc = strings.Replace(bc, "${profile}", pl.Profile, -1)
	bc = strings.Replace(bc, "${app}", pl.App, -1)
	ctx.ResponseBody(message, bc)
}

func (p *PipelineController) RemoteDeploy(pipelineService *service.PipelineService) error {
	remote, err := p.remoteService.InitRemote(pipelineService.Pipeline,"211855812371415399", pipelineService.BuildConfigs.Namespace, pipelineService.Namespace, pipelineService.App, pipelineService.Version)
	if err != nil {
		return nil
	}
	err = p.remoteService.Run(remote)
	return err
}
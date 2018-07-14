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
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"os"
	"strings"
	"github.com/hidevopsio/hicicd/pkg/service"
	"github.com/hidevopsio/hicicd/pkg/entity"
)

type CicdResponse struct {
	model.Response
}

// Operations about object
type PipelineController struct {
	BaseController
	PipelineService *service.PipelineService `inject:"pipelineService"`
	SelectorService *service.SelectorService `inject:"selectorService"`
}

const (
	ScmUrl      = "url"
	ScmUsername = "username"
	ScmPassword = "password"
	ScmToken    = "scmToken"
	ScmUid      = "uid"
)

func init() {
	web.Add(new(PipelineController))
}

func (c *PipelineController) Before(ctx *web.Context) {
	c.BaseController.Before(ctx)
}

// @Title Deploy
// @Description deploy application by the pipeline
// @Param	body
// @Success 200 {string}
// @Failure 403 body is empty
// @router / [post]
func (p *PipelineController) PostRun(ctx *web.Context) {
	log.Debug("CicdController.Run()")
	var pl entity.Pipeline
	err := ctx.RequestBody(&pl)
	// replace pl.Scm.Url with c.Url if it is empty
	if pl.Scm.Url == "" {
		pl.Scm.Url = p.Url
	}
	if err != nil {
		return
	}
	message := "success"
	if err == nil {
		selector, err := p.SelectorService.Get("1")
		if err != nil {
			return
		}
		selectorService := &service.PipelineService{}
		selectorService.Init(&pl, selector)
		go func() {
			err = selectorService.Run(p.Username, p.Password, p.ScmToken, p.Uid, false)
			if err != nil {
				message = err.Error()
			}
		}()
	} else {
		message = "failed, " + err.Error()
	}
	bc := os.Getenv("BUILD_CONSOLE")
	bc = strings.Replace(bc, "${project}", pl.Project, -1)
	bc = strings.Replace(bc, "${profile}", pl.Profile, -1)
	bc = strings.Replace(bc, "${app}", pl.App, -1)
	ctx.ResponseBody(message, bc)
}

func parseToken(claims jwt.MapClaims, prop string) string {
	return fmt.Sprintf("%v", claims[prop])
}

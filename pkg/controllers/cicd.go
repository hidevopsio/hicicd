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
	"github.com/hidevopsio/hicicd/pkg/ci/factories"
	"github.com/hidevopsio/hicicd/pkg/ci"
	"os"
	"strings"
)

type CicdResponse struct {
	model.Response
}

// Operations about object
type CicdController struct {
	BaseController
}

const (
	ScmUrl      = "url"
	ScmUsername = "username"
	ScmPassword = "password"
	ScmToken    = "scmToken"
	ScmUid      = "uid"
)

func init() {
	web.Add(new(CicdController))
}

func (c *CicdController) Before(ctx *web.Context) {
	c.BaseController.Before(ctx)
}

// @Title Deploy
// @Description deploy application by the pipeline
// @Param	body
// @Success 200 {string}
// @Failure 403 body is empty
// @router / [post]
func (c *CicdController) PostRun(ctx *web.Context) {
	log.Debug("CicdController.Run()")
	var pl ci.Pipeline
	err := ctx.RequestBody(&pl)
	// replace pl.Scm.Url with c.Url if it is empty
	if pl.Scm.Url == "" {
		pl.Scm.Url = c.Url
	}
	if err != nil {
		return
	}
	// invoke models
	pipelineFactory := new(factories.PipelineFactory)
	pipeline, err := pipelineFactory.New(pl.Name)
	message := "success"
	if err == nil {
		pipeline.Init(&pl)
		go func() {
			err = pipeline.Run(c.Username, c.Password, c.ScmToken, c.Uid, false)
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

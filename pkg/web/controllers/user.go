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
	"github.com/kataras/iris"
	"github.com/hidevopsio/hicicd/pkg/auth"
	"github.com/hidevopsio/hiboot/pkg/application"
	"time"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"os"
)

type UserRequest struct {
	Url      string `json:"url"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}


// Operations about object
type UserController struct {
}

// @Title Login
// @Description login
// @Param	body
// @Success 200 {string}
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Login(ctx iris.Context) {
	log.Debug("UserController.Login()")
	var request UserRequest

	err := ctx.ReadJSON(&request)
	if err != nil {
		application.ResponseError(ctx, err.Error(), iris.StatusInternalServerError)
		return
	}

	err = utils.Validate.Struct(&request)
	if err != nil {
		application.ResponseError(ctx, err.Error(), iris.StatusBadRequest)
		return
	}
	//log.Debug(token)
	var url string
	url = request.Url
	if url == "" {
		url = os.Getenv("SCM_URL")
	}
	if url == "" {
		application.ResponseError(ctx, err.Error(), iris.StatusInternalServerError)
	}else {
		// invoke models
		user := &auth.User{}
		_, message, err := user.Login(url, request.Username, request.Password)
		if err == nil {

				jwtToken, err := application.GenerateJwtToken(application.MapJwt{
					"url": url,
					"username": request.Username,
					"password": request.Password, // TODO: token is not working?
				}, 24, time.Hour)
				if err == nil {
					application.Response(ctx, message, &jwtToken)
				} else {
					application.ResponseError(ctx, err.Error(), iris.StatusInternalServerError)
				}

		} else {
			application.ResponseError(ctx, err.Error(), iris.StatusInternalServerError)
		}
	}
}

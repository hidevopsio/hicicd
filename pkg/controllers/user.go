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
	"github.com/hidevopsio/hicicd/pkg/auth"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"time"
	"github.com/hidevopsio/hiboot/pkg/log"
	"os"
	"net/http"
)

type UserRequest struct {
	Url      string `json:"url"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}


// Operations about object
type UserController struct {
	web.Controller
}

// new UserController instance
func init() {
	web.Add(new(UserController))
}

// @Title Login
// @Description login
// @Param	body
// @Success 200 {string}
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) PostLogin(ctx *web.Context) {
	log.Debug("UserController.Login()")
	var request UserRequest

	err := ctx.RequestBody(&request)
	if err != nil {
		return
	}

	url := request.Url
	if url == "" {
		url = os.Getenv("SCM_URL")
	}
	if url == "" {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
	}else {
		// invoke models
		user := &auth.User{}
		privateToken, uid, _, err := user.Login(url, request.Username, request.Password)
		if err == nil {

				jwtToken, err := web.GenerateJwtToken(web.JwtMap{
					"url": url,
					"username": request.Username,
					"password": request.Password, // TODO: token is not working?
					"scmToken": privateToken,
					"uid": uid,
				}, 24, time.Hour)
				if err == nil {
					ctx.ResponseBody("success", &jwtToken)
				} else {
					ctx.ResponseError(err.Error(), http.StatusInternalServerError)
				}

		} else {
			ctx.ResponseError(err.Error(), http.StatusForbidden)
		}
	}
}

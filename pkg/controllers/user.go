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
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"time"
	"github.com/hidevopsio/hiboot/pkg/log"
	"os"
	"net/http"
	"strconv"
	"github.com/hidevopsio/hiboot/pkg/starter/jwt"
)

type UserRequest struct {
	Url      string `json:"url"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Uid      string `json:"uid"`
}

// Operations about object
type UserController struct {
	web.Controller
	jwtToken jwt.Token
}

// new UserController instance
func init() {
	web.RestController(new(UserController))
}

func (c *UserController) Init(jwtToken jwt.Token) {
	c.jwtToken = jwtToken
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
	} else {
		// invoke models
		user := &auth.User{}
		privateToken, uid, _, err := user.Login(url, request.Username, request.Password)
		if err == nil {
			expired := os.Getenv("TOKEN_EXPIRED_HOUR")
			log.Debug("login expired time env:", expired)
			if expired == "" {
				expired = "24"
			}
			exp, err := strconv.ParseInt(expired, 10, 64)
			log.Debug("login expired time exp:", exp)
			jwtToken, err := c.jwtToken.Generate(jwt.Map{
				"url":      url,
				"username": request.Username,
				"password": request.Password, // TODO: token is not working?
				"scmToken": privateToken,
				"uid":      uid,
			}, 10, time.Minute)
			if err == nil {
				data := map[string]interface{}{
					"token": &jwtToken,
				}
				ctx.ResponseBody("success", &data)
			} else {
				ctx.ResponseError(err.Error(), http.StatusInternalServerError)
			}

		} else {
			ctx.ResponseError(err.Error(), http.StatusForbidden)
		}
	}
}

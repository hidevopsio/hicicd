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
	"github.com/dgrijalva/jwt-go"
	"strings"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
)

type BaseController struct {
	web.JwtController
	Username string
	Password string
	Url      string
	ScmToken string
}

func (c *BaseController) Before(ctx *web.Context) {
	log.Debug("route rule  add:{}")
	ti := ctx.Values().Get("jwt")
	token := ti.(*jwt.Token)
	var username, password string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Username = parseToken(claims, ScmUsername)
		c.Password = parseToken(claims, ScmPassword)
		c.ScmToken = parseToken(claims, ScmToken)
		log.Debugf("url: %v, username: %v, password: %v", username, strings.Repeat("*", len(password)))
	} else {
		log.Debug("valid username  password   err")

		return
	}
	ctx.Next()
}

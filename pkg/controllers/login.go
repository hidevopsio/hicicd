package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/prometheus/common/log"
	"net/url"
	"github.com/hidevopsio/hicicd/pkg/admin/auth"
	au "github.com/hidevopsio/hicicd/pkg/auth"
	"os"
	"net/http"
	"strconv"
	"time"
)

type LoginController struct {
	web.Controller
}


func init() {
	web.Add(new(LoginController))
}

func (l *LoginController) GetAuthUrl(ctx *web.Context)  {
	log.Debug("gitlab get oauth2 url")
	scmUrl := ctx.URLParam("url")
	if scmUrl == "" {
		scmUrl = os.Getenv("SCM_URL")
	}
	a := &auth.Auth{
		AuthURL: scmUrl,
		ApplicationId: auth.ApplicationId,
		CallbackUrl: auth.CallbackUrl,
	}
	url := a.GetAuthURL()
	responseBody := map[string]interface{}{
		"authUrl": &url,
	}
	ctx.ResponseBody("success", &responseBody)
}


func (l *LoginController) Get(ctx *web.Context) {
	log.Debug("gitlab oauth2 login ")
	code := ctx.URLParam("code")

	if code == "" {
		log.Error("loging get code not fount ")
		ctx.ResponseError("loging get code not fount ", http.StatusNotFound)
		return
	}
	//TODO 通过code  获取用户的accesstoken  和失效时间  通过时间来获取用户信息
	s := url.QueryEscape(auth.CallbackUrl)
	session := auth.NewClient(auth.BaseUrl, auth.AccessTokenUrl, auth.ApplicationId, s, auth.Secret)
	resp, err := session.GetAccessToken(code)
	if err != nil || resp.AccessToken == "" {
		log.Errorf("login session get accessToken error: %v", err)
		ctx.ResponseError("login session get accessToken error", http.StatusCreated)
		return
	}
	//TODO 通过accessToken 获取用户信息
	accessToken := resp.AccessToken
	log.Debugf("accessToken : %v", accessToken)
	lo := &au.Login{}
	user, err := lo.GetUser(os.Getenv("SCM_URL"), accessToken)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusCreated)
		return
	}
	log.Debug("oauth2 login user : %v ", user)
	ctx.ResponseBody("success", &user)
}

func (l *LoginController) PostAdmin(ctx *web.Context) {
	log.Debug("gitlab oauth2 login admin set usename and password")
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
		user := &au.User{}
		privateToken, uid, _, err := user.Login(url, request.Username, request.Password)
		if err == nil {
			expired := os.Getenv("TOKEN_EXPIRED_HOUR")
			log.Debug("login expired time env:", expired)
			if expired == "" {
				expired = "24"
			}
			exp, err:=strconv.ParseInt(expired, 10, 64)
			log.Debug("login expired time exp:", exp)
			jwtToken, err := web.GenerateJwtToken(web.JwtMap{
				"url": url,
				"username": request.Username,
				"password": request.Password, // TODO: token is not working?
				"scmToken": privateToken,
				"uid": uid,
			}, exp, time.Minute)
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
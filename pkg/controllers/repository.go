package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/info"
	"os"
	"strings"
	"github.com/hidevopsio/hicicd/pkg/orch/kong"
)

// Operations about object
type RepositoryController struct {
	BaseController
}

func init() {
	web.Add(new(RepositoryController))
}

func (c *RepositoryController) Before(ctx *web.Context) {
	c.BaseController.Before(ctx)
}

func (c *RepositoryController) PostAppType(ctx *web.Context) {
	log.Debug("Repository add:")
	var project scm.Project
	err := ctx.RequestBody(&project)
	if err != nil {
		return
	}
	t := new(info.TypeInfo)
	err = t.RepositoryType(c.Url, c.ScmToken, project.Ref, project.ID)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusPreconditionFailed)
		return
	}
	a := new(kong.ApiRequest)
	a.Name = project.Name + "-" + project.Namespace
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	baseUrl = strings.Replace(baseUrl, "${profile}", project.Profile, -1)
	api, err := a.Get(baseUrl)
	if err == nil {
		t.Uri = api.Uris[0]
	}else {
		uris := "/" + project.Namespace + "-" + project.Name
		t.Uri = strings.Replace(uris, "-", "/", -1)
	}
	host := os.Getenv("KONG_HOST")
	host = strings.Replace(host, "${profile}", project.Profile, -1)
	t.Host = host
	ctx.ResponseBody("success", &t)
	return
}

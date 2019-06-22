package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/info"
	"os"
	"strings"
	"github.com/kevholditch/gokong"
)

// Operations about object
type RepositoryController struct {
	BaseController
}

func init() {
	web.RestController(new(RepositoryController))
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
	err = t.RepositoryType(c.JwtProperty("url"), c.JwtProperty("scmToken"), project.Ref, project.ID)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusPreconditionFailed)
		return
	}
	name := project.Name + "-" + project.Namespace
	baseUrl := os.Getenv("KONG_ADMIN_URL")
	baseUrl = strings.Replace(baseUrl, "${profile}", project.Profile, -1)
	config := &gokong.Config{
		HostAddress: baseUrl,
	}
	api, err := gokong.NewClient(config).Apis().GetByName(name)
	if api != nil {
		t.Uri = api.Uris[0]
	} else {
		name := strings.Replace(project.Name, "-", "/", -1)
		t.Uri = "/" + project.Namespace + "/" + name
	}
	host := os.Getenv("KONG_HOST")
	host = strings.Replace(host, "${profile}", project.Profile, -1)
	t.Host = host
	ctx.ResponseBody("success", &t)
	return
}

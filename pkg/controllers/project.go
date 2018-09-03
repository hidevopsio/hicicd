package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"net/http"
	"github.com/hidevopsio/hicicd/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/log"
	"strconv"
)

// Operations about object
type ProjectController struct {
	BaseController
}

func init() {
	web.RestController(new(ProjectController))
}

func (c *ProjectController) Before(ctx *web.Context) {
	c.BaseController.Before(ctx)
}

func (c *ProjectController) Get(ctx *web.Context) {
	log.Debug("ProjectController All Projects()")
	search := ctx.URLParam("search")
	page, err := ctx.URLParamInt("page")
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	scmToken := c.JwtProperty("scmToken")
	url := c.JwtProperty("url")
	p := &app.Project{}
	projects, err := p.ListProjects(scmToken, url, search, page)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}

func (c *ProjectController) GetMember(ctx *web.Context) {
	log.Debug("ProjectController Project Member()")
	gid, err := ctx.URLParamInt("gid")
	pid, err := ctx.URLParamInt("pid")
	scmToken := c.JwtProperty("scmToken")
	url := c.JwtProperty("url")
	uid, _ := strconv.Atoi(c.JwtProperty("uid"))
	p := &app.Project{}
	projects, err := p.GetProjectMember(scmToken, url, pid, uid, gid)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}

package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"net/http"
	"github.com/hidevopsio/hicicd/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/log"
)

// Operations about object
type ProjectController struct {
	BaseController
}

func init() {
	web.Add(new(ProjectController))
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
	p := &app.Project{}
	projects, err := p.ListProjects(c.ScmToken, c.Url, search, page)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}

func (c *ProjectController) GetMember(ctx *web.Context)  {
	log.Debug("ProjectController Project Member()")
	gid, err := ctx.URLParamInt("gid")
	pid, err := ctx.URLParamInt("pid")
	p := &app.Project{}
	projects, err := p.GetProjectMember(c.ScmToken, c.Url, pid, c.Uid, gid)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}





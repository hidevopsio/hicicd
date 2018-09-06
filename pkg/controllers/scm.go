package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"net/http"
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/app"
)

// Operations about object
type ScmController struct {
	BaseController
}




func init() {
	web.Add(new(ScmController))
}

func (c *ScmController) Before(ctx *web.Context) {
	c.BaseController.Before(ctx)
}

func (c *ScmController) PostListGroups(ctx *web.Context) {
	log.Debug("ScmController.GetAllProject()")
	g := &app.Group{}
	groupMember, err := g.ListGroups(c.ScmToken, c.Url, c.Uid, 1)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &groupMember)
}

func (c *ScmController) PostListGroupProjects(ctx *web.Context) {
	log.Debug("ScmController.List Group Projects()")
	var g scm.Group
	err := ctx.RequestBody(&g)
	if err != nil {
		return
	}
	p := &app.Group{}
	projects, err := p.ListGroupProjects(c.ScmToken, c.Url, g.ID, g.Page)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}

func (c *ScmController) PostListProjects(ctx *web.Context) {
	log.Debug("ScmController All Projects()")
	var project scm.Project
	err := ctx.RequestBody(&project)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	p := &app.Project{}
	projects, err := p.ListProjects(c.ScmToken, c.Url, project.Search, project.Page)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}

func (c *ScmController) PostGetProjectMember(ctx *web.Context)  {
	log.Debug("ScmController Project Member()")
	var projectMember scm.ProjectMember
	err := ctx.RequestBody(&projectMember)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	p := &app.Project{}
	projects, err := p.GetProjectMember(c.ScmToken, c.Url, projectMember.Pid, c.Uid, projectMember.Gid)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}


func (c *ScmController) PostSearch(ctx *web.Context){
	log.Debug("ScmController Post Search()")
	var search scm.Search
	err := ctx.RequestBody(&search)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("ScmController ScmController search", search.Keyword)
	p := &app.Project{}
	projects, err := p.Search(c.Url, c.ScmToken, search.Keyword)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}

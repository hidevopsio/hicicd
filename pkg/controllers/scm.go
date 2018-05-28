package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/auth"
	"net/http"
	"github.com/hidevopsio/hicicd/pkg/scm"
)

type ProjectMember struct {
	Page      int    `json:"page"`
}
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
	permission := &auth.Permission{}
	groupMember, err := permission.ListGroups(c.ScmToken, c.Url, c.Uid)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &groupMember)
}

func (c *ScmController) PostListGroupProjects(ctx *web.Context) {
	log.Debug("ScmController.List Group Projects()")
	var group scm.Group
	err := ctx.RequestBody(&group)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	permission := &auth.Permission{}
	projects, err := permission.ListGroupProjects(c.ScmToken, c.Url, group.ID, group.Page)
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
	permission := &auth.Permission{}
	projects, err := permission.ListProjects(c.ScmToken, c.Url, project.Page)
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
	permission := &auth.Permission{}
	projects, err := permission.GetProjectMember(c.ScmToken, c.Url, projectMember.Pid, c.Uid)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}
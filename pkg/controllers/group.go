package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hicicd/pkg/app"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
)

// Operations about object
type GroupController struct {
	BaseController
}

func init() {
	web.Add(new(GroupController))
}

func (c *GroupController) Before(ctx *web.Context) {
	c.BaseController.Before(ctx)
}

func (c *GroupController) Get(ctx *web.Context) {
	log.Debug("GroupController.GetAllProject()")
	page, err := ctx.URLParamInt("page")
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	g := &app.Group{}
	groupMember, err := g.ListGroups(c.ScmToken, c.Url, c.Uid, page)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &groupMember)
}

func (c *GroupController) GetProjects(ctx *web.Context) {
	log.Debug("GroupController.List Group Projects()")
	page, err := ctx.URLParamInt("page")
	gid, err := ctx.URLParamInt("gid")
	g := &app.Group{}
	projects, err := g.ListGroupProjects(c.ScmToken, c.Url, gid, page)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}
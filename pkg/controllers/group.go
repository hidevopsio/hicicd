package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hicicd/pkg/app"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"strconv"
)

// Operations about object
type GroupController struct {
	BaseController
}

func init() {
	web.RestController(new(GroupController))
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
	scmToken := c.JwtProperty("scmToken")
	url := c.JwtProperty("url")
	uid, _ := strconv.Atoi(c.JwtProperty("uid"))
	g := &app.Group{}
	groupMember, err := g.ListGroups(scmToken, url, uid, page)
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
	scmToken := c.JwtProperty("scmToken")
	url := c.JwtProperty("url")
	g := &app.Group{}
	projects, err := g.ListGroupProjects(scmToken, url, gid, page)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	ctx.ResponseBody("success", &projects)
}

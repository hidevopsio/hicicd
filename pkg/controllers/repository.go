package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/info"
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
		ctx.ResponseError(err.Error(), http.StatusUnavailableForLegalReasons)
		return
	}
	t := new(info.TypeInfo)
	err = t.RepositoryType(c.Url, c.ScmToken, project.Ref, project.ID)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusPreconditionFailed)
		return
	}
	ctx.ResponseBody("success", &t)
	return
}


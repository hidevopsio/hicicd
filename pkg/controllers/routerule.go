package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/orch/istio"
	"net/http"
)



// Operations about object
type RouteruleController struct {
	BaseController
}

func init() {
	web.Add(new(RouteruleController))
}


func (c *RouteruleController) Before(ctx *web.Context) {
	c.BaseController.Before(ctx)
}

func (c *RouteruleController) PostAdd(ctx *web.Context) {
	log.Debug("route rule  add:{}")
	var rule  istio.RouterRule
	err := ctx.RequestBody(&rule)
	if err != nil {
		return
	}
	config, err := istio.NewClient()
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusRequestedRangeNotSatisfiable)
	}
	rule.Crd = config
	version, err := rule.Create()
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Debug("route rule create success get resource version:{}",version)
	r := IResponse{
		Version:version,
	}
	ctx.ResponseBody("success", &r)
	return
}

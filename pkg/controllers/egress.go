package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/orch/istio"
	"net/http"
)



// Operations about object
type EgressController struct {
	BaseController
}

func init() {
	web.Add(new(EgressController))
}


func (c *EgressController) Before(ctx *web.Context) {
	c.BaseController.Before(ctx)
}

func (c *EgressController) PostAdd(ctx *web.Context) {
	log.Debug("egress  add:{}")
	var egress  istio.Egress
	err := ctx.RequestBody(&egress)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusUnavailableForLegalReasons)
		return
	}
	config, err := istio.NewClient()
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusRequestedRangeNotSatisfiable)
	}
	egress.Crd = config
	version, err := egress.Create()
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Debug("egress create success get resource version:{}",version)
	r := IResponse{
		Version:version,
	}
	ctx.ResponseBody("success", &r)
	return
}

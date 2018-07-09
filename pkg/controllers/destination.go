package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/orch/istio"
	"net/http"
)



// Operations about object
type DestinationController struct {
	BaseController
}

func init() {
	web.Add(new(EgressController))
}


func (c *DestinationController) Before(ctx *web.Context) {
	c.BaseController.Before(ctx)
}

func (c *DestinationController) Post(ctx *web.Context) {
	log.Debug("destination  add:{}")
	var destination istio.Destination
	err := ctx.RequestBody(&destination)
	if err != nil {
		return
	}
	config, err := istio.NewClient()
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusRequestedRangeNotSatisfiable)
	}
	destination.Crd = config
	version, err := destination.Create()
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

package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"istio.io/istio/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/orch/istio"
	"net/http"
)

// Operations about object
type IstioController struct {
	BaseController
}
func init() {
	web.Add(new(IstioController))
}

type IResponse struct {
	Version string `json:"version"`
}


func (c *IstioController) PostDelete(ctx *web.Context)  {
	log.Debug("istio delete controller type")
	var route istio.RouterRule
	ctx.RequestBody(&route)
	if route.Type >= 5 {
		ctx.ResponseError("Index Out Of Bounds Exception route Type Length to Large", http.StatusLengthRequired)
		return
	}
	config, err := istio.NewClient()
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusRequestedRangeNotSatisfiable)
		return
	}
	route.Crd = config
	err = route.Delete(istio.Typ[route.Type])
	if err != nil {
		ctx.ResponseError( err.Error(), http.StatusRequestedRangeNotSatisfiable)
		return
	}
	message := "success"
	ctx.ResponseBody(message, nil)
}


func (c *IstioController) PostGet(ctx *web.Context)  {
	log.Debug("istio get controller type")

	var route istio.RouterRule
	ctx.RequestBody(&route)
	if route.Type >= 5 {
		ctx.ResponseError("Index Out Of Bounds Exception route Type Length to Large", http.StatusLengthRequired)
		return
	}
	config, err := istio.NewClient()
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusRequestedRangeNotSatisfiable)
		return
	}
	route.Crd = config
	conf, exists := route.Get(istio.Typ[route.Type])
	if exists == false {
		ctx.ResponseError( "not found istio type:" + istio.Typ[route.Type] , http.StatusNotFound)
		return
	}
	message := "success"
	ctx.ResponseBody(message, conf)
}
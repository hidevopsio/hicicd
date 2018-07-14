package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/entity"
	"github.com/hidevopsio/hicicd/pkg/service"
	"net/http"
)

type SelectorController struct {
	BaseController
	SelectorService *service.SelectorService `inject:"selectorService"`
}

func init() {
	web.Add(new(SelectorController))
}

func (s *SelectorController) Before(ctx *web.Context) {
	s.BaseController.Before(ctx)
}

func (s *SelectorController) Post(ctx *web.Context) {
	log.Debug("Selector add:")
	var selector entity.Selector
	err := ctx.RequestBody(&selector)
	if err != nil {
		return
	}
	err = s.SelectorService.Add(&selector)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusExpectationFailed)
	}
	ctx.ResponseBody("success", selector)
}


func (s *SelectorController) Get(ctx *web.Context) {
	log.Debug("Selector get")
	id := ctx.URLParam("id")
	selector, err := s.SelectorService.Get(id)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusExpectationFailed)
		return
	}
	ctx.ResponseBody("success", selector)
	return
}


func (s *SelectorController) Delete(ctx *web.Context) {
	log.Debug("Selector delete")
	id := ctx.URLParam("id")
	err := s.SelectorService.Delete(id)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusExpectationFailed)
		return
	}
	ctx.ResponseBody("success", nil)
	return
}
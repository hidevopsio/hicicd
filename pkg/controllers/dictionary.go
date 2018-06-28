package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hicicd/pkg/service"
	"github.com/hidevopsio/hicicd/pkg/admin"
	"net/http"
)

// Operations about object
type DictionaryController struct {
	BaseController
	DictionaryService *service.DictionaryService `component:"service"`
}

func init() {
	web.Add(new(DictionaryController))
}

func (d *DictionaryController) Post(ctx *web.Context) {
	dictionary := &admin.Dictionary{}
	err := ctx.RequestBody(dictionary)
	if err == nil {
		err := d.DictionaryService.CreateDictionaryService(dictionary)
		if err != nil {
			ctx.ResponseError(err.Error(), http.StatusExpectationFailed)
		}
		ctx.ResponseBody("success", dictionary)
	}

}

func (d *DictionaryController) Get(ctx *web.Context) {
	id := ctx.URLParam("id")
	dictionary, err := d.DictionaryService.GetDictionaryService(id)
	if err != nil {
		ctx.ResponseError("Resource is not found", http.StatusNotFound)
	} else {
		ctx.ResponseBody("success", dictionary)
	}
}

func (d *DictionaryController) Delete(ctx *web.Context) {
	id := ctx.URLParam("id")
	err := d.DictionaryService.DeleteDictionaryService(id)
	if err != nil {
		ctx.ResponseError("Resource is not found", http.StatusNotFound)
	} else {
		ctx.ResponseBody("success", nil)
	}
}


package controllers

import (
		"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hicicd/pkg/service"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/model"
)

// RestController
type DictionaryController struct {
	BaseController
	dictionaryService *service.DictionaryService
}

func init() {
	web.RestController(new(DictionaryController))
}

func (c *DictionaryController) Init(dictionaryService *service.DictionaryService) {
	c.dictionaryService = dictionaryService
}

func (c *DictionaryController) Before(ctx *web.Context) {
	log.Debug("controller before:{}")
	username := c.JwtProperty("username")
	log.Info(username)
	ctx.Next()
}

//Get 获取字典
func (c *DictionaryController) GetByType(id int32) (model.Response, error) {
	dictionary, err := c.dictionaryService.GetType(id)
	response := new(model.BaseResponse)
	response.SetData(dictionary)
	return response, err
}
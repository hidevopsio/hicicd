package controllers

import (
	"github.com/hidevopsio/hicicd/pkg/service"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/model"
)

// Operations about object
type ProfileController struct {
	BaseController
	profileService   *service.ProfileService
}

func (p *ProfileController) Init(profileService *service.ProfileService) {
	p.profileService = profileService
}

func init() {
	web.RestController(new(ProfileController))
}

func (p *ProfileController) Before(ctx *web.Context) {
	p.BaseController.Before(ctx)
}

func (p *ProfileController) GetById(id string) (model.Response, error) {
	profile, err := p.profileService.Get(id)
	response := new(model.BaseResponse)
	response.SetData(profile)
	return response, err
}
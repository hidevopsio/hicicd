package controllers

import "github.com/hidevopsio/hiboot/pkg/starter/web"

// Operations about object
type IstioController struct {
	web.JwtController
}

func init() {
	web.Add(new(IstioController))
}

func (c *IstioController) GetDelete(ctx *web.Context)  {
	
}

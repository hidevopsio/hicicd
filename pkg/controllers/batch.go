package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hicicd/pkg/entity"
	"github.com/hidevopsio/hicicd/pkg/service"
	"github.com/hidevopsio/hiboot/pkg/log"
	"net/http"
	"github.com/hidevopsio/hicicd/pkg/app"
	"github.com/hidevopsio/hicicd/pkg/info"
	"strings"
)

// Operations about object
type BatchController struct {
	BaseController
	SelectorService *service.SelectorService `inject:"selectorService"`
}

func init() {
	web.Add(new(BatchController))
}

func (p *BatchController) Before(ctx *web.Context) {
	p.BaseController.Before(ctx)
}

func (p *BatchController) Get(ctx *web.Context) {
	log.Debug("BatchController.Run()")
	search := ctx.URLParam("search")
	profile := ctx.URLParam("profile")
	pr := &app.Project{}
	projects, err := pr.ListProjects(p.ScmToken, p.Url, search, 1)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	for _, project := range projects  {
		t := new(info.TypeInfo)
		err = t.RepositoryType(p.Url, p.ScmToken, "development", project.ID)
		uri := project.Namespace + "-" + project.Name
		uri  = strings.Replace(uri, "-", "/", -1)
		pl := entity.Pipeline{
			App: project.Name,
			BuildConfigs: entity.BuildConfigs{
				Enable: true,
				Project: search,
			},
			Cluster: t.ClusterName,
			DeploymentConfigs: entity.DeploymentConfigs{
				ForceUpdate: true,
				Enable: true,
			},
			GatewayConfigs:entity.GatewayConfigs{
				Uri: uri,
				Enable: true,
			},
			IstioConfigs: entity.IstioConfigs{
				Enable: false,
			},
			Name: t.AppType,
			Profile: profile,
			Project: search,
			Scm: entity.Scm{
				Ref: "development",
			},
			Version: "v1",
		}
		if pl.Scm.Url == "" {
			pl.Scm.Url = p.Url
		}
		message := "success"
		selector, err := p.SelectorService.Get("1")
		if err != nil {
			return
		}
		selectorService := &service.PipelineService{}
		selectorService.Init(&pl, selector)
		go func() {
			err = selectorService.Run(p.Username, p.Password, p.ScmToken, p.Uid, false)
			if err != nil {
				message = err.Error()
			}
		}()
		ctx.ResponseBody(message, nil)
	}
}

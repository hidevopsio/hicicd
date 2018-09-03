package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/app/web"
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
	SelectorService *service.SelectorService
}

func init() {
	web.RestController(new(BatchController))
}

func (b *BatchController) Init(selectorService *service.SelectorService) {
	b.SelectorService = selectorService
}

func (b *BatchController) Before(ctx *web.Context) {
	b.BaseController.Before(ctx)
}

func (b *BatchController) Get(ctx *web.Context) {
	log.Debug("BatchController.Run()")
	search := ctx.URLParam("search")
	profile := ctx.URLParam("profile")
	pr := &app.Project{}
	projects, err := pr.ListProjects(b.ScmToken, b.Url, search, 1)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusInternalServerError)
		return
	}
	for _, project := range projects {
		t := new(info.TypeInfo)
		err = t.RepositoryType(b.Url, b.ScmToken, "development", project.ID)
		uri := project.Namespace + "-" + project.Name
		uri = strings.Replace(uri, "-", "/", -1)
		pl := entity.Pipeline{
			App: project.Name,
			BuildConfigs: entity.BuildConfigs{
				Enable:  true,
				Project: search,
			},
			Cluster: t.ClusterName,
			DeploymentConfigs: entity.DeploymentConfigs{
				ForceUpdate: true,
				Enable:      true,
			},
			GatewayConfigs: entity.GatewayConfigs{
				Uri:    uri,
				Enable: true,
			},
			IstioConfigs: entity.IstioConfigs{
				Enable: false,
			},
			Name:    t.AppType,
			Profile: profile,
			Project: search,
			Scm: entity.Scm{
				Ref: "development",
			},
			Version: "v1",
		}
		if pl.Scm.Url == "" {
			pl.Scm.Url = b.Url
		}
		message := "success"
		selector, err := b.SelectorService.Get("1")
		if err != nil {
			return
		}
		selectorService := &service.PipelineService{}
		selectorService.Initialize(&pl, selector)
		go func() {
			err = selectorService.Run(b.Username, b.Password, b.ScmToken, b.Uid, false)
			if err != nil {
				message = err.Error()
			}
		}()
		ctx.ResponseBody(message, nil)
	}
}

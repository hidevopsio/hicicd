package controllers

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/scm/factories"
	"encoding/base64"
	"gopkg.in/yaml.v2"
)

// Operations about object
type RepositoryController struct {
	BaseController
}

func init() {
	web.Add(new(RepositoryController))
}

//app.cluster-name
type Application struct {
	App    App `yaml:"app"`
	Spring Spring `yaml:"spring"`
}

type Spring struct {
	Profiles  Profiles `yaml:"profiles"`
}

type Profiles struct {
	Include []string `yaml:"include,flow"`
} 


type App struct {
	Project     string `yaml:"project"`
	ClusterName string `yaml:"cluster-name"`
}

type AppType struct {
	Name       string `json:"name"`
	FileName   string `json:"file_name"`
	FileType   string `json:"file_type"`
	configFile string `json:"config_file"`
}

type Body struct {
	ClusterName string `json:"cluster_name"`
	Name        string `json:"name"`
}

var AppTypeLUT = []AppType{
	{Name: "java", FileName: "pom.xml", FileType: "blob", configFile: "src/main/resources/application.yml"},
	{Name: "nodejs-dist", FileName: "dist", FileType: "tree",},
	{Name: "nodejs", FileName: "package.json", FileType: "blob"},
}

func (c *RepositoryController) Before(ctx *web.Context) {
	c.BaseController.Before(ctx)
}

func (c *RepositoryController) PostAppType(ctx *web.Context) {
	log.Debug("Repository add:{}")
	body := Body{}
	var project scm.Project
	err := ctx.RequestBody(&project)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusUnavailableForLegalReasons)
		return
	}
	scmFactory := new(factories.ScmFactory)
	r, err := scmFactory.NewResitory(factories.GitlabScmType)
	if err != nil {
		ctx.ResponseError(err.Error(), http.StatusLengthRequired)
		return
	}
	treeNodes, err := r.ListTree(c.Url, c.ScmToken, project.Ref, project.ID)
	for _, appType := range AppTypeLUT {
		for _, treeNode := range treeNodes {
			if appType.FileType == treeNode.Type && appType.FileName == treeNode.Name {
				log.Info(appType)
				message := "success"
				if appType.configFile != "" {
					context, err := r.GetRepository(c.Url, c.ScmToken, appType.configFile, project.Ref, project.ID)
					if err != nil {
						ctx.ResponseError(err.Error(), http.StatusLengthRequired)
						return
					}
					application, err := parse(context)
					if err != nil {
						ctx.ResponseError(err.Error(), http.StatusLengthRequired)
						return
					}
					body.Name = appType.Name

					for  _, include := range application.Spring.Profiles.Include{
						if include == "axon-jgroups" {
							body.ClusterName = application.App.ClusterName
							body.Name = "java-cqrs"
						}
					}
				} else {
					body.Name = appType.Name
					body.ClusterName = ""
				}
				ctx.ResponseBody(message, body)
				return
			}
		}
	}
	ctx.ResponseError("tree not found", http.StatusLengthRequired)
	return
}

func parse(content string) (*Application, error) {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}
	application := &Application{}
	err = yaml.Unmarshal([]byte(data), application)
	log.Info("application:{}", application)
	return application, nil
}

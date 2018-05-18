package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type ProductInterface interface {
	GetProject(token, baseUrl string, id interface{}) (*gitlab.Project, error)

	GetProjectLlist(token, baseUrl string) ([]*gitlab.Project, error)

	GetUserProject(token, baseUrl, name, namespace string) bool
}

type Product struct {
	Token     string `json:"token"`
	BaseUrl   string `json:"base_url"`
	ID        interface{}
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (p *Product) GetProject() (*gitlab.Project, error) {
	log.Debug("Product.GetProject()")
	c := gitlab.NewClient(&http.Client{}, p.Token)
	c.SetBaseURL(p.BaseUrl + ApiVersion)
	log.Debug("before c.Session.GetSession(so)")
	project, _, err := c.Projects.GetProject(p.ID)
	log.Debug("after c.Session.GetSession(so)")
	return project, err
}

func (p *Product) GetProjectLlist() ([]*gitlab.Project, error) {
	log.Debug("Product.GetProjectLlist()")
	log.Debugf("url: %v", p.BaseUrl)
	c := gitlab.NewClient(&http.Client{}, p.Token)
	c.SetBaseURL(p.BaseUrl + ApiVersion)
	project, _, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	log.Debugf("after project: %v", len(project))
	return project, err
}

func (p *Product) GetUserProject() bool {
	log.Debug("Product.GetUserProject()")
	log.Debugf("url: %v", p.BaseUrl)
	c := gitlab.NewClient(&http.Client{}, p.Token)
	c.SetBaseURL(p.BaseUrl + ApiVersion)
	log.Debug("before c.Session.GetProjectLlist")
	projects, _, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	if err != nil {
		log.Error("get list project :", err)
		return false
	}
	log.Debug("after c.Session.GetSession(so)")
	log.Debug("get project size: ", len(projects))
	for _, project := range projects {
		if project.Name == p.Name && project.Namespace.Name == p.Namespace {
			log.Debugf("project name: %v , name : %v", project.Name, p.Name)
			return true
		}
	}
	return false
}
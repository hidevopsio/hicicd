package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type Project struct {
	Token     string `json:"token"`
	BaseUrl   string `json:"base_url"`
	ID        interface{}
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (p *Project) GetProject() (*gitlab.Project, error) {
	log.Debug("Product.GetProject()")
	c := gitlab.NewClient(&http.Client{}, p.Token)
	c.SetBaseURL(p.BaseUrl + ApiVersion)
	log.Debug("before c.Session.GetSession(so)")
	project, _, err := c.Projects.GetProject(p.ID)
	log.Debug("after c.Session.GetSession(so)")
	return project, err
}

func (p *Project) GetProjectList() ([]*gitlab.Project, error) {
	log.Debug("Product.GetProjectList()")
	log.Debugf("url: %v", p.BaseUrl)
	c := gitlab.NewClient(&http.Client{}, p.Token)
	c.SetBaseURL(p.BaseUrl + ApiVersion)
	project, _, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	log.Debugf("after project: %v", len(project))
	return project, err
}

func (p *Project) GetUserProject(baseUrl, token, name, namespace string) (int, error) {
	log.Debug("Product.GetUserProject()")
	log.Debugf("url: %v", baseUrl)
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.Session.GetProjectList")
	projects, _, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	if err != nil {
		log.Error("get list project :", err)
		return 0, err
	}
	log.Debug("after c.Session.GetSession(so)")
	log.Debug("get project size: ", len(projects))
	for _, project := range projects {
		if project.Name == name && project.Namespace.Name == namespace {
			log.Debugf("project name: %v , name : %v", project.Name, name)
			return project.ID, nil
		}
	}
	return 0, err
}
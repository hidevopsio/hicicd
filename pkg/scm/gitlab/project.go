package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/jinzhu/copier"
)

type Project struct {
	scm.Project
}

func (p *Project) GetProject(baseUrl, id, token string) (int, int, error) {
	log.Debug("project.GetProject()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.project.GetProject(so)")
	project, _, err := c.Projects.GetProject(id)
	if err != nil {
		log.Error("Projects.GetProject err:", err)
		return 0, 0, err
	}
	return project.ID, project.Namespace.ID, err
}

func (p *Project) GetGroupId(url, token string, pid int) (int, error) {
	log.Debug("project.GetProject()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(p.BaseUrl + ApiVersion)
	log.Debug("before c.Session.GetSession(so)")
	project, _, err := c.Projects.GetProject(pid)
	log.Debug("after c.project.GetProject(so)", project)
	return project.Namespace.ID, err
}

func (p *Project) ListProjects(baseUrl, token, search string, page int) ([]scm.Project, error) {
	log.Debug("project.ListProjects()")
	log.Debugf("url: %v", baseUrl)
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	listProjectsOptions := &gitlab.ListProjectsOptions{}
	if search != "" {
		listProjectsOptions = &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
			Search: &search,
		}
	}else{
		listProjectsOptions = &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
		}
	}
	ps, _, err := c.Projects.ListProjects(listProjectsOptions)
	log.Debugf("after project: %v", len(ps))
	projects := []scm.Project{}
	project := &scm.Project{}
	for _, pro := range ps {
		copier.Copy(project, pro)
		projects = append(projects, *project)
	}
	return projects, err
}

func (p *Project) Search(baseUrl, token, search string) ([]scm.Project, error){
	log.Debug("Search.GetProjects()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before Search.project(so)", search)
	listProjectsOptions := &gitlab.ListProjectsOptions{
		Search: &search,
	}
	ps, _, err := c.Projects.ListProjects(listProjectsOptions)
	if err != nil {
		return nil, err
	}
	log.Debugf("after Search.project: %v", len(ps))
	projects := []scm.Project{}
	project := &scm.Project{}
	for _, pro := range ps {
		copier.Copy(project, pro)
		projects = append(projects, *project)
	}
	return projects, err
}
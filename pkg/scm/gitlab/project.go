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

func (p *Project) GetProject(baseUrl, name, token string) (*gitlab.Project, error) {
	log.Debug("project.GetProject()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.Session.GetSession(so)")
	project, _, err := c.Projects.GetProject(name)
	return project, err
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

func (p *Project) ListProjects(baseUrl, token string, page int) ([]scm.Project, error) {
	log.Debug("project.ListProjects()")
	log.Debugf("url: %v", baseUrl)
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	listProjectsOptions := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			Page: page,
		},
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

func (p *Project) ListUserProjects(baseUrl, token, name, namespace string) (int, error) {
	log.Debug("project.ListUserProjects()")
	log.Debugf("url: %v", baseUrl)
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.project.ListUserProjects")
	page := 1
	for {
		listProjectsOptions := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
		}
		projects, _, err := c.Projects.ListProjects(listProjectsOptions)
		if err != nil {
			log.Error("get list project :", err)
			return 0, err
		}
		if len(projects) == 0 {
			break
		}
		log.Debug("after c.project.project(so)")
		log.Debug("get project size: ", len(projects))
		for _, project := range projects {
			if project.Name == name && project.Namespace.Name == namespace {
				log.Debugf("project name: %v , name : %v", project.Name, name)
				return project.ID, nil
			}
		}
		page++
	}
	return 0, nil
}

func (p *Project) ListGroupProjects(baseUrl, token, namespace string) ([]scm.Project, error) {
	log.Debug("project List Group Project")
	c := gitlab.NewClient(&http.Client{}, token)
	scmProjects := []scm.Project{}
	scmProject := &scm.Project{}
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("c group projects ")
	page := 1
	for {
		listProjectsOptions := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		}
		projects, _, err := c.Projects.ListProjects(listProjectsOptions)
		if err != nil {
			return nil, err
		}
		if len(projects) == 0 {
			break
		}
		log.Debugf("project size : %v", len(projects))
		for _, p := range projects {
			if p.Namespace.Name == namespace {
				copier.Copy(scmProject, p)
				scmProjects = append(scmProjects, *scmProject)
			}
		}
		page++
	}
	return scmProjects, nil
}

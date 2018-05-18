package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type Project interface {


	GetProject(token, baseUrl string, id interface{}) (*gitlab.Project, error)

	GetProjectLlist(token, baseUrl string) ([]*gitlab.Project, error)

	GetUserProject(token, baseUrl, name,namespace string) bool
}




func (s *Session) GetProject(token, baseUrl string, id interface{}) (*gitlab.Project, error) {
	log.Debug("Session.GetSession()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.Session.GetSession(so)")
	project, _, err := c.Projects.GetProject(id)
	log.Debug("after c.Session.GetSession(so)")
	return project, err
}

func (s *Session) GetProjectLlist(token, baseUrl string) ([]*gitlab.Project, error) {
	log.Debug("Session.GetSession()")
	log.Debugf("url: %v", baseUrl)
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	project, _, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	log.Debugf("after project: %v", len(project) )
	return project, err
}

func (s *Session) GetUserProject(token, baseUrl, name,namespace string) bool{
	log.Debug("Session.GetSession()")
	log.Debugf("url: %v", baseUrl)
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.Session.GetProjectLlist")
	projects, _, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	if err != nil {
		log.Error("get list project :",err)
		return false
	}
	log.Debug("after c.Session.GetSession(so)")
	log.Debug("get project size: ", len(projects))
	for _,project := range projects {
		if project.Name == name && project.Namespace.Name == namespace{
			log.Debugf("project name: %v , name : %v", project.Name, name)
			return true
		}
	}
	return false
}

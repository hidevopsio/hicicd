package app

import (
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/scm/factories"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type Project struct {
	Project scm.ProjectInterface
	ProjectMember scm.ProjectMemberInterface
}

type ProjectInterface interface {
	ListProjects(token, baseUrl, search string, page int) ([]scm.Project,  error)
	GetProjectMember(token, baseUrl string, pid, uid, gid int) (int, error)
	Search(baseUrl, token, search string) ([]scm.Project, error)
}

func (p *Project) ListProjects(token, baseUrl, search string, page int) ([]scm.Project,  error)  {
	log.Debug("Permission List Group Projects:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	p.Project, err = scmFactory.NewProject(factories.GitlabScmType)
	if err != nil  {
		return nil, err
	}
	projects, err := p.Project.ListProjects(baseUrl, token, search, page)
	return projects, err
}

func (p *Project) GetProjectMember(token, baseUrl string, pid, uid, gid int) (int, error){
	log.Debug("Project List Group Projects:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	p.ProjectMember, err = scmFactory.NewProjectMember(factories.GitlabScmType)
	if err != nil  {
		return 0, err
	}
	projectMember, err := p.ProjectMember.GetProjectMember(token, baseUrl, pid, uid, gid)
	return projectMember.AccessLevelValue, err
}

func (p *Project) Search(baseUrl, token, search string) ([]scm.Project, error){
	log.Debug("permission project Search: {}")
	projects := []scm.Project{}
	scmFactory := new(factories.ScmFactory)
	var err error
	p.Project, err = scmFactory.NewProject(factories.GitlabScmType)
	if err != nil  {
		return projects, err
	}
	projects, err = p.Project.Search(baseUrl, token, search)
	return projects, err
}
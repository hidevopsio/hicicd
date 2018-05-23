package auth

import (
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/scm/factories"
)

type Permission struct {
	Project scm.ProjectInterface
	ProjectMember scm.ProjectMemberInterface
}

type PermissionInterface interface {
	Get(baseUrl, token, name, namespace string, uid int) (string, string, int, error)
}

func (p *Permission) Get(baseUrl, token, name, namespace string, uid int) (string, string, int, error) {
	scmFactory := new(factories.ScmFactory)
	var err error
	p.Project, err = scmFactory.NewProject(factories.GitlabScmType)
	if err != nil {
		return "", "", 0, err
	}
	pid, err := p.Project.GetUserProject(baseUrl, token, name, namespace)
	if err != nil {
		return "", "", 0, err
	}
	p.ProjectMember, err = scmFactory.NewProjectMember(factories.GitlabScmType)
	if err != nil {
		return "", "", 0, err
	}
	metaName, roleRefName, accessLevelValue, err := p.ProjectMember.GetProjectMember(token, baseUrl, pid, uid)
	return metaName, roleRefName, accessLevelValue, err
}
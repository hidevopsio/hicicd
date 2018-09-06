package auth

import (
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/scm/factories"
)

type Permission struct {
	Project scm.ProjectInterface
	ProjectMember scm.ProjectMemberInterface
	Group scm.GroupInterface
	GroupMember scm.GroupMemberInterface
}

const (
	NoPermissions        int = 0
	GuestPermissions     int = 10
	ReporterPermissions  int = 20
	DeveloperPermissions int = 30
	MasterPermissions    int = 40
	OwnerPermission      int = 50
)


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
	id := namespace + "/" + name
	pid, gid, err :=p.Project.GetProject(baseUrl, id, token)
	if err != nil {
		return "", "", 0, err
	}
	p.ProjectMember, err = scmFactory.NewProjectMember(factories.GitlabScmType)
	if err != nil {
		return "", "", 0, err
	}
	projectMember, err := p.ProjectMember.GetProjectMember(token, baseUrl, pid, uid, gid)
	return projectMember.MetaName, projectMember.RoleRefName, projectMember.AccessLevelValue, err
}

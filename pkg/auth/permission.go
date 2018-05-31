package auth

import (
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/scm/factories"
	"github.com/hidevopsio/hiboot/pkg/log"
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
	pid, err := p.Project.ListUserProjects(baseUrl, token, name, namespace)
	if err != nil {
		return "", "", 0, err
	}
	p.ProjectMember, err = scmFactory.NewProjectMember(factories.GitlabScmType)
	if err != nil {
		return "", "", 0, err
	}
	projectMember, err := p.ProjectMember.GetProjectMember(token, baseUrl, pid, uid)
	return projectMember.MetaName, projectMember.RoleRefName, projectMember.AccessLevelValue, err
}

func (p *Permission) ListGroups(token, baseUrl string, uid int) ([]scm.Group, error){
	log.Debug("Permission ListGroups:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	p.Group, err = scmFactory.NewGroup(factories.GitlabScmType)
	if err != nil  {
		return nil, err
	}
	groups, err := p.Group.ListGroups(token, baseUrl)
	scmGroups := []scm.Group{}
	for _, group := range groups {
		accessLevelValue, err := p.ListGroupMembers(token, baseUrl, group.ID, uid)
		if err != nil {
			return nil, err
		}
		group.AccessLevelValue = accessLevelValue
		scmGroups = append(scmGroups, group)
	}
	return scmGroups, err
}

func (p *Permission) GetGroup(token, baseUrl string, gid int) (*scm.Group, error){
	log.Debug("Permission GetGroup:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	p.Group, err = scmFactory.NewGroup(factories.GitlabScmType)
	if err != nil  {
		return nil, err
	}
	 scmGroup, err := p.Group.GetGroup(token, baseUrl, gid)
	return scmGroup, err
}

func (p *Permission) ListGroupMembers(token, baseUrl string, gid, uid int) (int,  error)  {
	log.Debug("Permission ListGroupMembers:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	p.GroupMember, err = scmFactory.NewGroupMember(factories.GitlabScmType)
	if err != nil  {
		return 0, err
	}
	accessLevelValue, err := p.GroupMember.ListGroupMembers(token, baseUrl, gid, uid)
	return accessLevelValue, err
}

func (p *Permission) ListGroupProjects(token, baseUrl string, gid, page int) ([]scm.Project,  error)  {
	log.Debug("Permission List Group Projects:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	p.Group, err = scmFactory.NewGroup(factories.GitlabScmType)
	if err != nil  {
		return nil, err
	}
	projects, err := p.Group.ListGroupProjects(token, baseUrl, gid, page)
	return projects, err
}

func (p *Permission) ListProjects(token, baseUrl string, page int) ([]scm.Project,  error)  {
	log.Debug("Permission List Group Projects:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	p.Project, err = scmFactory.NewProject(factories.GitlabScmType)
	if err != nil  {
		return nil, err
	}
	projects, err := p.Project.ListProjects(baseUrl, token, page)
	return projects, err
}

func (p *Permission) GetProjectMember(token, baseUrl string, pid, uid int) (int, error){
	log.Debug("Permission List Group Projects:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	p.ProjectMember, err = scmFactory.NewProjectMember(factories.GitlabScmType)
	if err != nil  {
		return 0, err
	}
	projectMember, err := p.ProjectMember.GetProjectMember(token, baseUrl, pid, uid)
	return projectMember.AccessLevelValue, err
}
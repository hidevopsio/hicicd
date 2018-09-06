package app

import (
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/scm/factories"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type Group struct {
	Group scm.GroupInterface
	GroupMember scm.GroupMemberInterface
}

type GroupInterface interface {
	ListGroups(token, baseUrl string, uid int) ([]scm.Group, error)
	GetGroup(token, baseUrl string, gid int) (*scm.Group, error)
	ListGroupProjects(token, baseUrl string, gid, page int) ([]scm.Project,  error)
	ListGroupMembers(token, baseUrl string, gid, uid int) (int,  error)
}

func (g *Group) ListGroups(token, baseUrl string, uid, page int) ([]scm.Group, error){
	log.Debug("Group ListGroups:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	g.Group, err = scmFactory.NewGroup(factories.GitlabScmType)
	if err != nil  {
		return nil, err
	}
	groups, err := g.Group.ListGroups(token, baseUrl, page)
	scmGroups := []scm.Group{}
	for _, group := range groups {
		accessLevelValue, err := g.ListGroupMembers(token, baseUrl, group.ID, uid)
		if err != nil {
			return nil, err
		}
		group.AccessLevelValue = accessLevelValue
		scmGroups = append(scmGroups, group)
	}
	return scmGroups, err
}

func (g *Group) GetGroup(token, baseUrl string, gid int) (*scm.Group, error){
	log.Debug("Group GetGroup:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	g.Group, err = scmFactory.NewGroup(factories.GitlabScmType)
	if err != nil  {
		return nil, err
	}
	scmGroup, err := g.Group.GetGroup(token, baseUrl, gid)
	return scmGroup, err
}

func (g *Group) ListGroupMembers(token, baseUrl string, gid, uid int) (int,  error)  {
	log.Debug("Permission ListGroupMembers:{}")
	scmFactory := new(factories.ScmFactory)
	var err error
	g.GroupMember, err = scmFactory.NewGroupMember(factories.GitlabScmType)
	if err != nil  {
		return 0, err
	}
	accessLevelValue, err := g.GroupMember.ListGroupMembers(token, baseUrl, gid, uid)
	return accessLevelValue, err
}

func (p *Group) ListGroupProjects(token, baseUrl string, gid, page int) ([]scm.Project,  error)  {
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
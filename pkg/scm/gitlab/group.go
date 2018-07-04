package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/jinzhu/copier"
)

type Group struct {
	scm.Group
}

func (g *Group) ListGroups(token, baseUrl string, page int) ([]scm.Group, error) {
	log.Debug("group.ListGroups()")
	scmGroups := []scm.Group{}
	scmGroup := &scm.Group{}
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.group.ListGroups(so)")
	opt := &gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{
			Page: page,
		},
	}
	groups, _, err := c.Groups.ListGroups(opt)
	if err != nil {
		return nil, err
	}
	log.Debug("after c.Group.ListGroups(so)")
	for _, group := range groups {
		copier.Copy(scmGroup, group)
		scmGroups = append(scmGroups, *scmGroup)
	}
	return scmGroups, err
}

func (g *Group) GetGroup(token, baseUrl string, gid int) (*scm.Group, error) {
	log.Debug("group.GetGroup()")
	scmGroup := &scm.Group{}
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.group.ListGroups(so)")
	group, _, err := c.Groups.GetGroup(gid)
	log.Debug("after c.Session.GetSession(so)")
	if err != nil {
		return nil, err
	}
	copier.Copy(scmGroup, group)
	return scmGroup, err
}

func (g *Group) GetGroupMembers(token, baseUrl string, uid int) ([]scm.Group, error) {
	log.Debug("group.ListGroups()")
	scmGroups := []scm.Group{}
	scmGroup := &scm.Group{}
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.group.ListGroups(so)")
	opt := &gitlab.ListGroupsOptions{}
	groups, _, err := c.Groups.ListGroups(opt)
	if err != nil {
		return nil, err
	}
	log.Debug("after c.Group.ListGroups(so)")
	for _, group := range groups {
		copier.Copy(scmGroup, group)
		gid := group.ID
		groupMembers, _, err := c.Groups.ListGroupMembers(gid, &gitlab.ListGroupMembersOptions{})
		if err != nil {
			return nil, err
		}
		for _, groupMember := range groupMembers {
			if groupMember.ID == uid {
				for id, permissions := range scm.Permissions {
					if groupMember.AccessLevel == id {
						scmGroup.AccessLevelValue = permissions.AccessLevelValue
						break
					}
				}
			}
		}
		scmGroups = append(scmGroups, *scmGroup)
	}
	return scmGroups, err
}

func (g *Group) ListGroupProjects(token, baseUrl string, gid, page int) ([]scm.Project, error) {
	log.Debug("group.ListGroups()")
	scmProjects := []scm.Project{}
	scmProject := &scm.Project{}
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.group.ListGroups(so)")
	opt := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			Page: page,
		},
	}
	projects, _, err := c.Groups.ListGroupProjects(gid, opt)
	log.Debug("ListGroupProjects :{}",len(projects))
	if err != nil {
		log.Error("Group ListGroupProjects :{}", err)
		return nil, err
	}
	for _, project := range projects {
		copier.Copy(scmProject, project)
		scmProjects = append(scmProjects, *scmProject)
	}
	return scmProjects, err
}
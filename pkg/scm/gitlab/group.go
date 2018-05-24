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



func (g *Group) ListGroups(token, baseUrl string) ([]*scm.Group, error)  {
	log.Debug("group.ListGroups()")
	scmGroups := []*scm.Group{}
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
	for group := range groups {
		copier.Copy(scmGroup, group)
		scmGroups = append(scmGroups, scmGroup)
	}
	return scmGroups, err

}

func (g *Group) GetGroup(token, baseUrl string, gid int) error {
	log.Debug("group.GetGroup()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.group.ListGroups(so)")
	group, _, err := c.Groups.GetGroup(gid)
	log.Debug("after c.Session.GetSession(so)")
	if err != nil {
		return err
	}
	copier.Copy(g, group)
	return nil
}
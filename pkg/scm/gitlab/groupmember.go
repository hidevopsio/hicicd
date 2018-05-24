package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm"
	copier "github.com/jinzhu/copier"
)

type GroupMember struct {
	scm.GroupMember
}

func (gm *GroupMember) ListGroupMembers(token, baseUrl string, gid int) ([]*scm.GroupMember,  error)  {
	log.Debug("group.ListGroups()")
	scmGroupMembers := []*scm.GroupMember{}
	scmGroupMember := &scm.GroupMember{}
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.group.ListGroups(so)")
	opt := &gitlab.ListGroupMembersOptions{}
	groupMembers, _, err := c.Groups.ListGroupMembers(gid, opt)
	if err != nil {
		return nil, err
	}
	log.Debug("after c.Session.GetSession(so)")
	for groupMember := range groupMembers{
		copier.Copy(scmGroupMember, groupMember)
		scmGroupMembers = append(scmGroupMembers, scmGroupMember)
	}
	return scmGroupMembers, err

}

func (gm *GroupMember) GetGroupMember(token, baseUrl string, gid, uid int) error  {
	log.Debug("group.ListGroups()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.group.ListGroups(so)")
	opt := &gitlab.ListGroupMembersOptions{}
	groupMembers, _, err := c.Groups.ListGroupMembers(gid, opt)
	log.Debug("after c.group member.groupMembers(so)")
	if err != nil {
		return err
	}
	for _, groupMember := range groupMembers {
		if groupMember.ID == uid {
			copier.Copy(gm, groupMember)
			return nil
		}
	}
	return nil
}
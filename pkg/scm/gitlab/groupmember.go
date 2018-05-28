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


func (gm *GroupMember) GetGroupMember(token, baseUrl string, gid, uid int) (*scm.GroupMember,  error)  {
	log.Debug("group.ListGroups()")
	scmGroupMember := &scm.GroupMember{}
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.group.ListGroups(so)")
	opt := &gitlab.ListGroupMembersOptions{}
	groupMembers, _, err := c.Groups.ListGroupMembers(gid, opt)
	log.Debug("after c.group member.groupMembers(so)")
	if err != nil {
		return nil, err
	}
	for _, groupMember := range groupMembers {
		if groupMember.ID == uid {
			copier.Copy(scmGroupMember, groupMember)
		}
	}
	return scmGroupMember, nil
}


func (gm *GroupMember) ListGroupMembers(token, baseUrl string, gid, uid int) (int,  error)  {
	log.Debug("group.ListGroups()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.group.ListGroups(so)")
	opt := &gitlab.ListGroupMembersOptions{}
	groupMembers, _, err := c.Groups.ListGroupMembers(gid, opt)
	if err != nil {
		return 0, err
	}
	log.Debug("after gm.GroupMember.ListGroupMembers(so)")
	for _, groupMember := range groupMembers{
		if groupMember.ID == uid {
			for id, permissions := range scm.Permissions  {
				if groupMember.AccessLevel == id {
					return permissions.AccessLevelValue, nil
				}
			}
		}
	}
	return 0, err

}
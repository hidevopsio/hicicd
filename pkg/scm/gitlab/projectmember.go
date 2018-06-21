package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm"
)
type ProjectMember struct {
	scm.ProjectMember
}



func (p *ProjectMember) GetProjectMember(token, baseUrl string, pid, uid, gid int) (scm.ProjectMember, error) {
	log.Debug("Product.GetProject()")
	scmProjectMember := scm.ProjectMember{}
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before p.project.GetProjectMember(so)", pid)
	opt := &gitlab.ListGroupMembersOptions{
		ListOptions: gitlab.ListOptions{
			Page: 50,
		},
	}
	groupMembers, _, err := c.Groups.ListGroupMembers(gid, opt)
	if err != nil {
		log.Error("Groups.ListGroupMembers err:", err)
		return scmProjectMember, err
	}
	log.Debug("Groups.ListGroupMembers ")
	for _, groupMember := range groupMembers {
		if groupMember.ID == uid {
			for id, permissions := range scm.Permissions  {
				if groupMember.AccessLevel == id {
					scmProjectMember.ProjectPermissions = permissions
					return scmProjectMember, nil
				}
			}
		}
	}
	projectMember, _, err := c.Projects.GetProjectMember(pid, uid)
	if err != nil {
		log.Error("ProjectMembers.GetProjectMember ", err)
		return scmProjectMember, err
	}
	log.Debug("after c.Session.GetSession(so)")
	for id, permissions := range scm.Permissions  {
		if projectMember.AccessLevel == id {
			scmProjectMember.ProjectPermissions = permissions
			return scmProjectMember, nil
		}
	}
	return scmProjectMember, err
}

func (p *ProjectMember) ListProjectMembers(token, baseUrl string, pid int)  ([]*gitlab.ProjectMember, error) {
	log.Debug("Product.GetProject()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.Session.GetSession(so)")
	opt := &gitlab.ListProjectMembersOptions{}
	projectMembers, _, err := c.Projects.ListProjectMembers(pid, opt)
	if err != nil {
		return nil, err
	}
	log.Debug("after c.Session.GetSession(so)")

	return projectMembers, nil
}
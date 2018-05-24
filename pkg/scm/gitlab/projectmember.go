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



func (p *ProjectMember) GetProjectMember(token, baseUrl string, pid, uid int) (string, string, int, error) {
	log.Debug("Product.GetProject()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.Session.GetSession(so)")
	projectMember, _, err := c.ProjectMembers.GetProjectMember(pid, uid)
	if err != nil {
		return "", "", 0, err
	}
	log.Debug("after c.Session.GetSession(so)")
	for id, permissions := range scm.Permissions  {
		if projectMember.AccessLevel == id {
			return permissions.MetaName, permissions.RoleRefName, permissions.AccessLevelValue, nil
		}
	}
	return "", "", 0, err
}

func (p *ProjectMember) ListProjectMembers(token, baseUrl string, pid int)  ([]*gitlab.ProjectMember, error) {
	log.Debug("Product.GetProject()")
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.Session.GetSession(so)")
	opt := &gitlab.ListProjectMembersOptions{}
	projectMembers, _, err := c.ProjectMembers.ListProjectMembers(pid, opt)
	if err != nil {
		return nil, err
	}
	log.Debug("after c.Session.GetSession(so)")

	return projectMembers, nil
}
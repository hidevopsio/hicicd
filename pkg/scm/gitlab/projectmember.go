package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type ProjectPermissions struct {
	MetaName    string `json:"meta_name"`
	RoleRefName string `json:"role_ref_name"`
}

var Permissions = map[gitlab.AccessLevelValue]ProjectPermissions{
	gitlab.OwnerPermission:      ProjectPermissions{MetaName: "admin", RoleRefName: "admin"},
	gitlab.MasterPermissions:    ProjectPermissions{MetaName: "admin", RoleRefName: "admin"},
	gitlab.DeveloperPermissions: ProjectPermissions{MetaName: "edit-hptg8", RoleRefName: "edit"},
	gitlab.ReporterPermissions:  ProjectPermissions{MetaName: "view-gbtpw", RoleRefName: "view"},
	gitlab.GuestPermissions:     ProjectPermissions{MetaName: "view-gbtpw", RoleRefName: "view"},
	gitlab.NoPermissions:        ProjectPermissions{MetaName: "view-gbtpw", RoleRefName: "view"},
}

type ProjectMember struct {
	Token     string `json:"token"`
	BaseUrl   string `json:"base_url"`
	User      int    `json:"user"`
	Pid       int    `json:"pid"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (p *ProjectMember) GetProjectMember() (*gitlab.ProjectMember, error) {
	log.Debug("Product.GetProject()")
	c := gitlab.NewClient(&http.Client{}, p.Token)
	c.SetBaseURL(p.BaseUrl + ApiVersion)
	log.Debug("before c.Session.GetSession(so)")
	projectMember, _, err := c.ProjectMembers.GetProjectMember(p.Pid, p.User)
	log.Debug("after c.Session.GetSession(so)")
	return projectMember, err
}

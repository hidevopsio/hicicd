package scm

import (
	"github.com/xanzy/go-gitlab"
)

type ProjectMemberInterface interface {
	GetProjectMember(token, baseUrl string, pid, uid int) (string, string, int, error)
}


type ProjectPermissions struct {
	MetaName         string `json:"meta_name"`
	RoleRefName      string `json:"role_ref_name"`
	AccessLevelValue int    `json:"access_level_value"`
}

var Permissions = map[gitlab.AccessLevelValue]ProjectPermissions{
	gitlab.OwnerPermission:      ProjectPermissions{MetaName: "admin",      RoleRefName: "admin", AccessLevelValue: 50},
	gitlab.MasterPermissions:    ProjectPermissions{MetaName: "admin",      RoleRefName: "admin", AccessLevelValue: 40},
	gitlab.DeveloperPermissions: ProjectPermissions{MetaName: "edit-hptg8", RoleRefName: "edit",  AccessLevelValue: 30},
	gitlab.ReporterPermissions:  ProjectPermissions{MetaName: "view-gbtpw", RoleRefName: "view",  AccessLevelValue: 20},
	gitlab.GuestPermissions:     ProjectPermissions{MetaName: "view-gbtpw", RoleRefName: "view",  AccessLevelValue: 10},
	gitlab.NoPermissions:        ProjectPermissions{MetaName: "view-gbtpw", RoleRefName: "view",  AccessLevelValue: 0},
}

type ProjectMember struct {
	Token     string `json:"token"`
	BaseUrl   string `json:"base_url"`
	User      int    `json:"user"`
	Pid       int    `json:"pid" `
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
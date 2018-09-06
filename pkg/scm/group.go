package scm

type GroupInterface interface {
	ListGroups(token, baseUrl string, page int) ([]Group, error)
	GetGroup(token, baseUrl string, gid int) (*Group, error)
	ListGroupProjects(token, baseUrl string, gid, page int) ([]Project, error)
}

type Group struct {
	ID                   int        `json:"id"`
	Name                 string     `json:"name"`
	Path                 string     `json:"path"`
	Description          string     `json:"description"`
	LFSEnabled           bool       `json:"lfs_enabled"`
	AvatarURL            string     `json:"avatar_url"`
	WebURL               string     `json:"web_url"`
	RequestAccessEnabled bool       `json:"request_access_enabled"`
	FullName             string     `json:"full_name"`
	FullPath             string     `json:"full_path"`
	ParentID             int        `json:"parent_id"`
	Projects             []*Project `json:"projects"`
	AccessLevelValue     int        `json:"access_level_value"`
	Page                 int        `json:"page"`
}

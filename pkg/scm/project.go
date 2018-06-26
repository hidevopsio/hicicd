package scm

import "time"

type ProjectInterface interface {
	ListProjects(baseUrl, token, search string, page int) ([]Project, error)
	GetGroupId(url, token string, pid int) (int, error)
	GetProject(baseUrl, id, token string) (int, int, error)
	Search(baseUrl, token, search string) ([]Project, error)
}

type Project struct {
	ID                                        int        `json:"id"`
	Description                               string     `json:"description"`
	DefaultBranch                             string     `json:"default_branch"`
	Public                                    bool       `json:"public"`
	SSHURLToRepo                              string     `json:"ssh_url_to_repo"`
	HTTPURLToRepo                             string     `json:"http_url_to_repo"`
	WebURL                                    string     `json:"web_url"`
	TagList                                   []string   `json:"tag_list"`
	Name                                      string     `json:"name"`
	NameWithNamespace                         string     `json:"name_with_namespace"`
	Path                                      string     `json:"path"`
	PathWithNamespace                         string     `json:"path_with_namespace"`
	IssuesEnabled                             bool       `json:"issues_enabled"`
	OpenIssuesCount                           int        `json:"open_issues_count"`
	MergeRequestsEnabled                      bool       `json:"merge_requests_enabled"`
	ApprovalsBeforeMerge                      int        `json:"approvals_before_merge"`
	JobsEnabled                               bool       `json:"jobs_enabled"`
	WikiEnabled                               bool       `json:"wiki_enabled"`
	SnippetsEnabled                           bool       `json:"snippets_enabled"`
	ContainerRegistryEnabled                  bool       `json:"container_registry_enabled"`
	CreatedAt                                 *time.Time `json:"created_at,omitempty"`
	LastActivityAt                            *time.Time `json:"last_activity_at,omitempty"`
	CreatorID                                 int        `json:"creator_id"`
	ImportStatus                              string     `json:"import_status"`
	ImportError                               string     `json:"import_error"`
	Archived                                  bool       `json:"archived"`
	AvatarURL                                 string     `json:"avatar_url"`
	SharedRunnersEnabled                      bool       `json:"shared_runners_enabled"`
	ForksCount                                int        `json:"forks_count"`
	StarCount                                 int        `json:"star_count"`
	RunnersToken                              string     `json:"runners_token"`
	PublicJobs                                bool       `json:"public_jobs"`
	OnlyAllowMergeIfPipelineSucceeds          bool       `json:"only_allow_merge_if_pipeline_succeeds"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool       `json:"only_allow_merge_if_all_discussions_are_resolved"`
	LFSEnabled                                bool       `json:"lfs_enabled"`
	RequestAccessEnabled                      bool       `json:"request_access_enabled"`
	CIConfigPath                              *string    `json:"ci_config_path"`
	AccessLevel                               int        `json:"access_level"`
	Token                                     string     `json:"token"`
	BaseUrl                                   string     `json:"base_url"`
	Namespace                                 string     `json:"namespace"`
	Page                                      int        `json:"page"`
	Search                                    string     `json:"search"`
	Group                                     *Group
	Ref                                       string     `json:"ref"`
	Profile                                   string     `json:"profile"`
}

type Search struct {
	Keyword string `json:"keyword"`
}

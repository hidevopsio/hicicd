package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type ProjectMember struct {
	Token     string `json:"token"`
	BaseUrl   string `json:"base_url"`
	User      int
	Pid       int
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

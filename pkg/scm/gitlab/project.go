package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/jinzhu/copier"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func (s *Session) GetProject(baseUrl, username, password string) error {
	log.Debug("Session.GetSession()")
	log.Debugf("url: %v, username: %v", baseUrl, username)
	so := &gitlab.GetSessionOptions{
		Login:    &username,
		Password: &password,
	}
	c := gitlab.NewClient(&http.Client{}, "")
	c.SetBaseURL(baseUrl + ApiVersion)
	log.Debug("before c.Session.GetSession(so)")
	project, _, err := c.Projects.GetProject(so)
	log.Debug("after c.Session.GetSession(so)")

	copier.Copy(s, project)

	return err
}

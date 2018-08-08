package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/jinzhu/copier"
)

type User struct {
	scm.User
}


func (s *User) GetUser(baseUrl, accessToken string) (*scm.User, error) {
	log.Debug("Session get user")
	c := gitlab.NewOAuthClient(&http.Client{}, accessToken)
	c.SetBaseURL(baseUrl + ApiVersion)
	user, _, err := c.Users.CurrentUser()
	if err != nil {
		return nil, err
	}
	scmUser := &scm.User{}
	copier.Copy(scmUser, user)
	log.Debugf("get User : %v", user)
	return scmUser, nil

}

package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"time"
)

type User struct {
	ID               int             `json:"id"`
	Username         string          `json:"username"`
	Email            string          `json:"email"`
	Name             string          `json:"name"`
	State            string          `json:"state"`
	CreatedAt        *time.Time      `json:"created_at"`
	Bio              string          `json:"bio"`
	Skype            string          `json:"skype"`
	Linkedin         string          `json:"linkedin"`
	Twitter          string          `json:"twitter"`
	WebsiteURL       string          `json:"website_url"`
	ExternUID        string          `json:"extern_uid"`
	Provider         string          `json:"provider"`
	ThemeID          int             `json:"theme_id"`
	ColorSchemeID    int             `json:"color_scheme_id"`
	IsAdmin          bool            `json:"is_admin"`
	AvatarURL        string          `json:"avatar_url"`
	CanCreateGroup   bool            `json:"can_create_group"`
	CanCreateProject bool            `json:"can_create_project"`
	ProjectsLimit    int             `json:"projects_limit"`
	CurrentSignInAt  *time.Time      `json:"current_sign_in_at"`
	LastSignInAt     *time.Time      `json:"last_sign_in_at"`
	TwoFactorEnabled bool            `json:"two_factor_enabled"`
}

func (s *User) GetUser(baseUrl, accessToken string) error {
	log.Debug("Session get user")
	c := gitlab.NewOAuthClient(&http.Client{}, accessToken)
	c.SetBaseURL(baseUrl + ApiVersion)
	user, _, err := c.Users.CurrentUser()
	if err != nil {
		return err
	}
	log.Debugf("get User : %v", user)
	return nil

}

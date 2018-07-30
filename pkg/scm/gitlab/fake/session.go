package fake

import (
	"github.com/xanzy/go-gitlab"
	"errors"
)

type FakeSession struct {
	client *Client
}

type GetSessionOptions struct {
	Login    *string `url:"login,omitempty" json:"login,omitempty"`
	Email    *string `url:"email,omitempty" json:"email,omitempty"`
	Password *string `url:"password,omitempty" json:"password,omitempty"`
}

func TestGetSessionOptions() GetSessionOptions {
	login := "test"
	email := "test.com"
	password := "test"
	options := GetSessionOptions{
		Login:    &login,
		Email:    &email,
		Password: &password,
	}
	return options
}

func (s *FakeSession) GetSession(opt *gitlab.GetSessionOptions, options ... gitlab.OptionFunc) (*gitlab.Session, *gitlab.Response, error) {
	sessionOptions := TestGetSessionOptions()
	if opt.Login == sessionOptions.Login && opt.Password == sessionOptions.Password {
		session := &gitlab.Session{
			ID:           190,
			PrivateToken: "KASDUIO2901JSADJLK",
			Username : "test",
			Name: "test",
		}
		return session, nil, nil
	}

	return nil, nil, errors.New("Account password does not exist ")
}

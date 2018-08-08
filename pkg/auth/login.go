package auth

import (
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/scm/factories"
)

type Login struct {
	user scm.UserInterface
}

type LoginInterface interface {
	getUser(baseUrl, accessToken string) (*scm.User, error)
}



func (l *Login) GetUser(baseUrl, accessToken string) (*scm.User, error){
	scmFactory := new(factories.ScmFactory)
	var err error
	l.user, err = scmFactory.NewUser(factories.GitlabScmType)
	if err != nil {
		return nil, err
	}
	user, err := l.user.GetUser(baseUrl, accessToken)
	return user, err
}
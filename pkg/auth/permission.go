package auth

import (
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/hidevopsio/hicicd/pkg/scm/factories"
)

type Permission struct {
	Product scm.ProjectInterface
	ProductMember scm.ProjectMemberInterface
}

type PermissionInterface interface {
	Get(baseUrl, token, name, namespace string, uid int) (string, string, int, error)
}

func (a *Permission) Get(baseUrl, token, name, namespace string, uid int) (string, string, int, error) {
	scmFactory := new(factories.ScmFactory)
	var err error
	a.Product, err = scmFactory.NewProject(factories.GitlabScmType)
	if err != nil {
		return "", "", 0, err
	}
	pid, err := a.Product.GetUserProject(baseUrl, token, name, namespace)
	if err != nil {
		return "", "", 0, err
	}
	a.ProductMember, err = scmFactory.NewProjectMember(factories.GitlabScmType)
	if err != nil {
		return "", "", 0, err
	}
	metaName, roleRefName, accessLevelValue, err := a.ProductMember.GetProjectMember(token, baseUrl, pid, uid)
	return metaName, roleRefName, accessLevelValue, err
}
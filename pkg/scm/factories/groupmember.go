package factories

import (
	"github.com/hidevopsio/hicicd/pkg/scm"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/log"
	"errors"
)

func (s *ScmFactory) NewGroupMember(provider int) (scm.GroupMemberInterface, error)  {
	log.Debug("scm.NewSession()")
	switch provider {
	case GithubScmType:
		return nil, errors.New(fmt.Sprintf("SCM of type %d not implemented\n", provider))
	case GitlabScmType:
		return new(gitlab.GroupMember), nil
	default:
		return nil, errors.New(fmt.Sprintf("SCM of type %d not recognized\n", provider))
	}
	return nil, nil
}

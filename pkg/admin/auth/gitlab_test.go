package auth

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"strings"
	"github.com/prometheus/common/log"
)

func TestAuthUrl(t *testing.T)  {
	a := &Auth{
		AuthURL: "gitlab.vpclub:8022",
		ApplicationId: ApplicationId,
		CallbackUrl: CallbackUrl,
	}
	baseUrl := a.GetAuthURL()
	log.Infof("baseUrl: %v", baseUrl)
	assert.Equal(t, true, strings.Contains(baseUrl, "gitlab"))
}

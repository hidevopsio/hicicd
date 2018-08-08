package gitlab

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestUser_GetUser(t *testing.T) {
	token := "0767a31c454f00c9a7682871732b4ddb94be611a87717ae12e3b8db2842c3146"
	baseUrl := os.Getenv("SCM_URL")
	log.Debugf("accessToken: %v", token)
	user := new(User)
	u, err := user.GetUser(baseUrl, token)
	assert.Equal(t, nil, err)
	assert.Equal(t, "chulei", u.Username)
}

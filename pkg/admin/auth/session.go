package auth

import (
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/juju/errors"
)

const code = "09fa0bbe85fe580a0c4ab9f07b8da37e63724a053783ab17e8686ca9761c13bc"

type Session struct {
	AuthURL       string
	ApplicationId string
	Secret        string
	TokenURL      string
	profileURL    string
	CallbackUrl   string
	Code          string
}

type SessionRespons struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	CreateAt         string `json:"create_at"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewClient(authUrl, tokenUrl, applicationId, callbackUrl, secret string) *Session {
	s := &Session{
		AuthURL:       authUrl,
		TokenURL:      tokenUrl,
		ApplicationId: applicationId,
		CallbackUrl:   callbackUrl,
		Secret:        secret,
	}
	return s
}

func (session *Session) GetAccessToken(code string) (*SessionRespons, error) {
	log.Info("session GetAccessToken code : ", code)
	if code == "" {
		return nil, errors.New("code not nil")
	}
	session.Code = code
	t := utils.GetMatches(AccessTokenUrl)
	baseUrl := ReplaceEnv(AccessTokenUrl, t, session)
	sessionRespons := &SessionRespons{}
	_, err := Client("POST", baseUrl, sessionRespons)
	log.Info(sessionRespons)
	return sessionRespons, err
}

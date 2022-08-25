package wechatwork

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"water-reminder/pkg/httpclient"
	"water-reminder/pkg/wechatwork/constant"
)

type Application interface {
	SendMessage(msg *Message) error
}

type application struct {
	companyId string
	secret    string
	agentId   int

	token AccessToken
}

func (app *application) SendMessage(msg *Message) error {
	token, err := app.token.Token()
	if err != nil {
		return errors.Wrap(err, "get token")
	}
	msg.AgentId = app.agentId

	log.Debugf("SendMessage: %+v, token: %s", msg, token)
	resp, err := httpclient.Post(fmt.Sprintf(constant.URLSendMessage, token), msg)
	if err != nil {
		return errors.Wrap(err, "httpclient.Post")
	}
	return IsError(resp)
}

var _ Application = (*application)(nil)

func NewApplication(companyId, secret string, agentId int) (Application, error) {
	log.Debugf("NewApplication companyId: %s, secret: %s", companyId, secret)
	if len(companyId) == 0 {
		return nil, errors.New("企业 ID 不能为空")
	}
	if agentId == 0 {
		return nil, errors.New("应用 AgentId 不能为空")
	}
	if len(secret) == 0 {
		return nil, errors.New("应用 Secret 不能为空")
	}
	app := &application{companyId: companyId, agentId: agentId, secret: secret}
	token, err := NewAccessToken(app)
	if err != nil {
		return nil, errors.Wrap(err, "NewAccessToken")
	}
	app.token = token
	return app, nil
}

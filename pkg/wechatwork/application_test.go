package wechatwork

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

const (
	companyId = "ww48d2720d9851f7af"
	agentId   = 1000002
	secret    = "SYOdciKqJlKSkNXfjpeifKPa0NsKhBuq4oyq7to1wbI"
)

var app Application

func TestNewApplication(t *testing.T) {
	var err error
	app, err = NewApplication(companyId, secret, agentId)
	if err != nil {
		t.Fatal(errors.Wrap(err, "NewApplication"))
	}
}

func TestApplication_SendMessage(t *testing.T) {
	if app == nil {
		t.FailNow()
	}
	if err := app.SendMessage(NewTextMessage("test message")); err != nil {
		t.Error(err)
	}
}

package app

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
	"water-reminder/config"
	"water-reminder/pkg/wechatwork"
)

func Run(cfg *config.Config) {
	app, err := wechatwork.NewApplication(cfg.CompanyId, cfg.Secret, cfg.AgentId)
	if err != nil {
		log.Fatal(err)
	}

	var local *time.Location
	if local, err = time.LoadLocation(cfg.Location); err != nil {
		log.Fatal(errors.Wrapf(err, "LoadLocation %s", cfg.Location))
	}

	schedule, err := NewSchedule("09:00-18:00", "11:30-13:00", time.Hour*2, 3000, local, app)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(schedule.Start())
}

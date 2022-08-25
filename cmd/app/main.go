package main

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
	"water-reminder/config"
	"water-reminder/internal/app"
)

func init() {
	// 设置日志打印格式
	format := new(log.TextFormatter)
	format.FullTimestamp = true
	format.TimestampFormat = "06-01-02 15:04:05"
	log.SetFormatter(format)

	// 设置随机种子
	rand.Seed(time.Now().UnixNano())
}

func main() {
	cfg := config.NewConfig()
	if len(cfg.Location) == 0 {
		cfg.Location = "Asia/Shanghai"
	}
	log.Infof("NewConfig: %+v", cfg)
	app.Run(cfg)
}

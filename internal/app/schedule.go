package app

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
	"water-reminder/pkg/wechatwork"
)

type Schedule struct {
	// 一天时间的开启和结束：09:00-18:00
	oneDay *period
	// 休息时间：11:30-13:00
	breakTime *period

	interval time.Duration

	location *time.Location

	app wechatwork.Application

	target *ScheduleTarget
}

type ScheduleTarget struct {
	Target  int
	each    int
	current int

	// 已提醒提醒次数
	times int

	lastRemind map[int]bool
}

func fmtML(ml int) string {
	if ml < 1000 {
		return fmt.Sprintf("%dml", ml)
	}
	return fmt.Sprintf("%.1fL", float64(ml)/1000)
}

func (st *ScheduleTarget) reset() {
	st.each = 600
	st.current = 0
	st.times = 0
	st.lastRemind = map[int]bool{}
}

func (st *ScheduleTarget) message() string {
	switch st.times {
	case 1:
		var powerList = []string{
			"工作再忙也不要忘记喝水呐，听到没！",
			"坚持就是胜利💪，请收下这碗鸡汤😜",
			"PS 今天也是想念我宝贝的一天哦😘",
			"本提醒是您专属的哦，独一份！👏",
			"再忍一忍，周末马上就到了 😄",
			"我只是个无情的机器人，给我喝水👊",
			"人工智能的时代何时到来，这样我就可以下班啦😌",
			"这才是v0.1版本哦，一整个期待住了🤩",
		}
		power := powerList[rand.Int()%len(powerList)]
		return fmt.Sprintf(`今天是 %s ☀️, 憨憨来提醒你喝水啦！
今日喝水目标：%s
每次需要喝 %s 哦
%s`, time.Now().Format("06-01-02"), fmtML(st.Target), fmtML(st.each), power)
	default:
		return fmt.Sprintf(`叮咚，该喝水💦啦～
这是今天第 %d 次了昂
今日喝水目标 (%.0f%%) %s/%s🎯`, st.times, float64(st.current)/float64(st.Target)*100, fmtML(st.current), fmtML(st.Target))
	}
}

func (st *ScheduleTarget) remind(times int) string {
	switch times {
	case 0:
		return fmt.Sprintf(`我来啦，刚才提醒你喝水，喝够了没！
当前喝水目标 %s/%s🎯
喝够回复1️⃣，不够快给我喝去！`, fmtML(st.current), fmtML(st.Target))
	case 1:
		return fmt.Sprintf(`我又来啦，别嫌我烦，刚才给你说的目标喝够了没？
当前喝水目标 %s/%s🎯`, fmtML(st.current), fmtML(st.Target))
	default:
		return "什么也没有，可能系统出错了！"
	}
}

const periodLayout = "15:04"

// 时间段
type period struct {
	start time.Time
	end   time.Time
}

func (p *period) String() string {
	return fmt.Sprintf("%s-%s", p.start.Format(periodLayout), p.end.Format(periodLayout))
}

func (p *period) toTime(t time.Time) (time.Time, time.Time) {
	start := time.Date(t.Year(), t.Month(), t.Day(), p.start.Hour(), p.start.Minute(), 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), p.end.Hour(), p.end.Minute(), 0, 0, t.Location())
	return start, end
}

func parsePeriod(t string, local *time.Location) (*period, error) {
	const sep = "-"
	if len(t) == 0 {
		return nil, errors.New("parameter is empty")
	}
	sp := strings.Split(t, sep)
	if len(sp) != 2 {
		return nil, errors.Errorf("%s strings.Split result is not 2", t)
	}

	var err error
	var period = new(period)
	if period.start, err = time.ParseInLocation(periodLayout, sp[0], local); err != nil {
		return nil, err
	}
	if period.end, err = time.ParseInLocation(periodLayout, sp[1], local); err != nil {
		return nil, err
	}
	return period, nil
}

func NewSchedule(oneDay, breakTime string, interval time.Duration, target int, location *time.Location, app wechatwork.Application) (*Schedule, error) {
	if location == nil {
		location = time.Local
	}
	oneDayPeriod, err := parsePeriod(oneDay, location)
	if err != nil {
		return nil, err
	}
	breakTimePeriod, err := parsePeriod(breakTime, location)
	if err != nil {
		return nil, err
	}
	return &Schedule{oneDay: oneDayPeriod, breakTime: breakTimePeriod, interval: interval, location: location, app: app, target: &ScheduleTarget{Target: target}}, nil
}

func (s *Schedule) Start() error {
	log.Infof("Schedule.Start oneDay: %v, breakTime: %v, interval: %v", s.oneDay, s.breakTime, s.interval)
	for {
		s.target.reset()
		log.Infof("Schedule.Start %+v", s.target)
		next := s.delay()
		log.Infof("Schedule.Start next day after %v", next)
		<-time.After(next)
	}
}

func (s *Schedule) delay() time.Duration {
	nowTime := time.Now().In(s.location)
	log.Infof("Schedule.delay now time: %v", nowTime)

	startTime, endTime := s.oneDay.toTime(nowTime)
	log.Infof("Schedule.delay day start: %v", startTime)
	log.Infof("Schedule.delay day end: %v", endTime)
	breakStart, breakEnd := s.breakTime.toTime(nowTime)
	log.Infof("Schedule.delay break start: %v", breakStart)
	log.Infof("Schedule.delay break end: %v", breakEnd)

	breakDiff := breakEnd.Sub(breakStart)
	for {
		nextStart := startTime.Add(s.interval)
		if nextStart.After(breakStart) && nextStart.Before(breakEnd) {
			// 下次时间正好在休息时间
			if breakDiff >= s.interval {
				nextStart = breakEnd
			} else {
				nextStart = breakEnd.Add(s.interval - (breakStart.Sub(startTime)))
			}
		}

		// 判断当前时间是否已在计划中
		if nowTime.After(startTime) {
			startTime = nextStart
			continue
		}

		if startTime.After(endTime) {
			break
		}
		log.Infof("Schedule.delay schedule %v", startTime)
		s.schedule(startTime.Sub(nowTime))
		if startTime.Equal(endTime) {
			break
		}
		startTime = nextStart
	}

	// 零点时间
	todayTime := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, nowTime.Location())
	// 今天已经过的时间，一天的总时间减去已经过去的时间为下次触发循环的时间
	passed := nowTime.Sub(todayTime)
	log.Infof("Schedule.delay today passed %v", passed)
	const oneDay = 86400 * time.Second
	return oneDay - passed
}

func (s *Schedule) schedule(duration time.Duration) {
	log.Infof("Schedule.schedule %v", duration)
	time.AfterFunc(duration, func() {
		if s.target.current >= s.target.Target {
			return
		}
		s.target.times++
		s.target.current += s.target.each
		if err := s.app.SendMessage(wechatwork.NewTextMessage(s.target.message())); err != nil {
			log.Errorf("SendMessage %v", err)
		}
	})

	// 两次额外提醒，提示喝水量需达标
	// 1/4时提醒一次
	log.Info("Schedule.schedule remind0: ", duration+s.interval/4)
	time.AfterFunc(duration+s.interval/4, func() {
		if s.target.lastRemind[0] {
			return
		}
		if err := s.app.SendMessage(wechatwork.NewTextMessage(s.target.remind(0))); err != nil {
			log.Errorf("SendMessage %v", err)
		}
		if s.target.current >= s.target.Target {
			s.target.lastRemind[0] = true
		}
	})
	// 3/5时提醒一次
	log.Info("Schedule.schedule remind1: ", duration+s.interval*3/5)
	time.AfterFunc(duration+s.interval*3/4, func() {
		if s.target.lastRemind[1] {
			return
		}
		if err := s.app.SendMessage(wechatwork.NewTextMessage(s.target.remind(1))); err != nil {
			log.Errorf("SendMessage %v", err)
		}
		if s.target.current >= s.target.Target {
			s.target.lastRemind[1] = true
		}
	})
}

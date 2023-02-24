package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"sync"
	"time"
)

type Cron struct {
	mu        sync.RWMutex
	movements dsd.AutoMovementSlice

	crontab *cron.Cron
	infos   [][]ScheduleInfo
	infoCh  chan ScheduleInfo
}

type ScheduleInfo struct {
	CronID     int
	Function   dsd.CronFunction
	FuncID     int
	AutoHoming int
	start      time.Time
	end        time.Time
}

func New(movements dsd.AutoMovementSlice) *Cron {
	c := &Cron{
		movements: movements,
		crontab:   cron.New(cron.WithSeconds()),
		infos:     make([][]ScheduleInfo, 7),
		infoCh:    make(chan ScheduleInfo, 1),
	}

	if err := c.convert(movements); err != nil {
		log.Errorf(err.Error())
		c.infos = make([][]ScheduleInfo, 7)
	}

	return c
}

func (c *Cron) Start() <-chan ScheduleInfo {
	c.crontab.AddJob("@every 3s", c)
	c.crontab.Start()

	return c.infoCh
}

func (c *Cron) Stop() {
	c.crontab.Stop()
}

func (c *Cron) Run() {
	weekday := time.Now().Weekday()
	weekdayInfo := c.infos[weekday]

	now := time.Now()

	for _, info := range weekdayInfo {
		if now.After(info.start) && now.Before(info.end) {
			c.infoCh <- info
			return
		}
	}
}

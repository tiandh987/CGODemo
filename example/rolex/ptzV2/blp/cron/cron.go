package cron

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
	"sort"
	"sync"
	"time"
)

// Function 定时功能
type Function int

const (
	None       Function = iota // None
	Preset                     // 预置点
	Cruise                     // 巡航
	Trace                      // 巡迹
	LineScan                   // 线性扫描
	RegionScan                 // 区域扫描
)

func (f Function) Validate() error {
	if f < None || f > RegionScan {
		return errors.New("invalid cron function")
	}

	return nil
}

type Cron struct {
	mu        sync.RWMutex
	movements []dsd.PtzAutoMovement
	infos     [][]ScheduleInfo

	crontab *cron.Cron
	infoCh  chan ScheduleInfo
	quitCh  chan struct{}
}

type ScheduleInfo struct {
	CronID     int
	Function   Function
	FuncID     int
	start      time.Time
	end        time.Time
	autoHoming time.Duration
	sendTime   time.Time
}

func New(movements []dsd.PtzAutoMovement) (*Cron, error) {
	cron := &Cron{
		crontab: cron.New(cron.WithSeconds()),
		infos:   make([][]ScheduleInfo, 7),
		infoCh:  make(chan ScheduleInfo, 1),
		quitCh:  make(chan struct{}, 1),
	}

	if err := cron.convert(movements); err != nil {
		return nil, err
	}

	return cron, nil
}

func (c *Cron) List() []dsd.PtzAutoMovement {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.movements
}

func (c *Cron) Default() error {

	return nil
}

func (c *Cron) Set(movement *dsd.PtzAutoMovement) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	movements := c.movements
	for index, m := range movements {
		if m.ID == movement.ID {
			movements[index] = *movement
			break
		}
	}

	cron := &Cron{}
	if err := cron.convert(movements); err != nil {
		return err
	}

	c.movements = cron.movements
	c.infos = cron.infos

	return nil
}

func (c *Cron) Start() {
	c.crontab.AddJob("@every 3s", c)
	c.crontab.Start()
}

func (c *Cron) Stop() {
	c.crontab.Stop()
}

func (c *Cron) Quit() {
	ctx := c.crontab.Stop()
	select {
	case <-ctx.Done():
	}

	c.quitCh <- struct{}{}
	close(c.quitCh)
	close(c.infoCh)
}

func (c *Cron) InfoCh() <-chan ScheduleInfo {
	return c.infoCh
}

func (c *Cron) QuitCh() <-chan struct{} {
	return c.quitCh
}

func (c *Cron) Run() {
	weekday := time.Now().Weekday()
	weekdayInfo := c.infos[weekday]

	now := time.Now()
	if now.Before(weekdayInfo[0].start) || now.After(weekdayInfo[len(weekdayInfo)-1].end) {
		return
	}

	for index, info := range weekdayInfo {
		if now.After(info.start) && now.Before(info.end) {
			if now.Sub(info.sendTime) < info.autoHoming {
				return
			}
			c.infos[weekday][index].sendTime = now
			c.infoCh <- info
		}
	}
}

func (c *Cron) convert(movements []dsd.PtzAutoMovement) error {
	for _, movement := range movements {
		if !movement.Enable {
			continue
		}

		for weekday, sections := range movement.Schedule.WeekDay {
			for i, s := range sections.Section {
				start, err := parseTimeStr(s.TimeStr[0])
				if err != nil {
					return errors.New(fmt.Sprintf("cron id[%d] weekday [%d] section[%d %s] parse start failed", movement.ID, weekday, i))
				}

				end, err := parseTimeStr(s.TimeStr[1])
				if err != nil {
					return errors.New(fmt.Sprintf("cron id[%d] weekday [%d] section[%d %s] parse end failed", movement.ID, weekday, i))
				}

				if start.After(end) {
					return errors.New(fmt.Sprintf("cron id[%d] weekday [%d] section[%d] start must less than end", movement.ID, weekday, i))
				}

				funcID := 0
				switch Function(movement.Function) {
				case None:
					continue
				case Preset:
					funcID = movement.PresetID
				case Cruise:
					funcID = movement.TourID
				case Trace:
					funcID = movement.PatternID
				case LineScan:
					funcID = movement.LinearScanID
				case RegionScan:
					funcID = movement.RegionScanID
				default:
					return errors.New("invalid cron function")
				}

				info := ScheduleInfo{
					CronID:     int(movement.ID),
					Function:   Function(movement.Function),
					FuncID:     funcID,
					start:      start,
					end:        end,
					autoHoming: time.Duration(movement.AutoHoming.Time),
				}

				c.infos[weekday] = append(c.infos[weekday], info)
			}
		}
	}

	for weekday, info := range c.infos {
		sort.Slice(info, func(i, j int) bool {
			return info[i].start.Before(info[j].start)
		})

		for i := 0; i < len(info)-1; i++ {
			if info[i].end.After(info[i+1].start) {
				log.Errorf("invalid schedule section.\n%+v\n%+v", info[i], info[i+1])
				return errors.New(fmt.Sprintf("invalid schedule section, id[%d %d] weekday[%d]", info[i].CronID, info[i+1].CronID, weekday))
			}
		}
	}

	return nil
}

func parseTimeStr(str string) (time.Time, error) {
	timeStr := fmt.Sprintf("%s-%s-%s %s", time.Now().Format("2006"), time.Now().Format("01"),
		time.Now().Format("02"), str)

	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		log.Errorf(err.Error())
		return time.Time{}, err
	}

	return t, nil
}

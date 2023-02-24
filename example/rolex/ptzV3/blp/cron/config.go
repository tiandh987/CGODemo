package cron

import (
	"errors"
	"fmt"
	"github.com/tiandh987/CGODemo/example/rolex/config"
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"sort"
	"time"
)

func (c *Cron) List() dsd.AutoMovementSlice {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.movements
}

func (c *Cron) Default() error {
	//slice := dsd.NewAutoMovementSlice()
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

	cron := &Cron{
		infos: make([][]ScheduleInfo, 7),
	}
	if err := cron.convert(movements); err != nil {
		return err
	}

	if err := config.SetConfig(c.movements.ConfigKey(), c.movements); err != nil {
		return err
	}

	c.movements = movements
	c.infos = cron.infos

	return nil
}

func (c *Cron) convert(movements []dsd.PtzAutoMovement) error {
	for _, movement := range movements {
		if !movement.Enable {
			continue
		}

		for weekday, sections := range movement.Schedule.WeekDay {
			for i, s := range sections.Section {
				// 解析开始时间
				start, err := parseTimeStr(s.TimeStr[0])
				if err != nil {
					log.Error(err.Error())

					retErr := fmt.Errorf("cron id[%d] weekday [%d] section[%d] parse start failed",
						movement.ID, weekday, i)
					return retErr
				}

				// 解析结束时间
				end, err := parseTimeStr(s.TimeStr[1])
				if err != nil {
					log.Error(err.Error())

					retErr := fmt.Errorf("cron id[%d] weekday [%d] section[%d] parse end failed",
						movement.ID, weekday, i)
					return retErr
				}

				// 开始时间必须小于结束时间
				if start.After(end) {
					retErr := fmt.Errorf("cron id[%d] weekday [%d] section[%d] start must less than end",
						movement.ID, weekday, i)
					return retErr
				}

				// 获取 Function、FuncID
				funcID, err := movement.GetFuncID()
				if err != nil {
					return err
				}

				// 自动归位最小为 3s
				homing := movement.AutoHoming.Time
				if homing < 3 {
					homing = 3
				}

				info := ScheduleInfo{
					CronID:     int(movement.ID),
					Function:   dsd.CronFunction(movement.Function),
					FuncID:     funcID,
					AutoHoming: homing,
					start:      start,
					end:        end,
				}

				c.infos[weekday] = append(c.infos[weekday], info)
			}
		}
	}

	if err := c.checkSchedule(); err != nil {
		return err
	}

	c.printScheduleInfo()

	return nil
}

func (c *Cron) printScheduleInfo() {
	for weekday, info := range c.infos {
		log.Infof("[weekday:%d]\n", weekday)
		for _, i := range info {
			log.Infof("CronID:%d Function:%d FuncID:%d start:%s end:%s homing:%d\n", i.CronID, i.Function,
				i.FuncID, i.start.Format("15:04:05"), i.end.Format("15:04:05"), i.AutoHoming)
		}
	}
}

func (c *Cron) checkSchedule() error {
	for weekday, info := range c.infos {
		sort.Slice(info, func(i, j int) bool {
			return info[i].start.Before(info[j].start)
		})

		for i := 0; i < len(info)-1; i++ {
			if info[i].end.After(info[i+1].start) {
				log.Errorf("invalid schedule section (weekday:%d)\n"+
					"CronID:%d Function:%d FuncID:%d start:%s end:%s\n"+
					"CronID:%d Function:%d FuncID:%d start:%s end:%s\n",
					weekday,
					info[i].CronID, info[i].Function, info[i].FuncID, info[i].start.Format("15:04:05"),
					info[i].end.Format("15:04:05"),
					info[i+1].CronID, info[i+1].Function, info[i+1].FuncID, info[i+1].start.Format("15:04:05"),
					info[i+1].end.Format("15:04:05"))

				return errors.New(fmt.Sprintf("invalid schedule section"))
			}
		}
	}

	return nil
}

func parseTimeStr(str string) (time.Time, error) {
	timeStr := fmt.Sprintf("%s-%s-%s %s", time.Now().Format("2006"), time.Now().Format("01"),
		time.Now().Format("02"), str)

	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id        int64
	Cfg       string
	Executer  string
	Cron      string
	Name      string
	CanelFunc func() error
	Ctime     int64
	Utime     int64
}

var parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom |
	cron.Month | cron.Dow | cron.Descriptor)

func (j Job) NextTime() time.Time {
	s, _ := parser.Parse(j.Cron)
	return s.Next(time.Now())
}

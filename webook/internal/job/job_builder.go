package job

import "github.com/robfig/cron/v3"

type CronJobBuilder struct {
}

func (b *CronJobBuilder) Build(job Job) cron.Job {

}

type CronJobAdapter func() error

func (c *CronJobAdapter) Run() {

}

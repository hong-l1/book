package job

import (
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"time"
)

type CronJobBuilder struct {
	p *prometheus.SummaryVec
	l logger.Loggerv1
}

func NewCronJobBuilder(l logger.Loggerv1) CronJobBuilder {
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "gobook",
		Subsystem: "webook",
		Name:      "cron_job",
	}, []string{"name"})
	prometheus.MustRegister(summary)
	return CronJobBuilder{
		l: l,
		p: summary,
	}
}
func (b *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return CronJobAdapter(func() error {
		start := time.Now()
		err := job.Run()
		if err != nil {
			b.l.Error("任务执行失败",
				logger.String("name", name),
				logger.Error(err))
		}
		defer func() {
			duration := time.Since(start)
			b.p.WithLabelValues(name).Observe(duration.Seconds())
		}()
		return nil
	})
}

type CronJobAdapter func() error

func (c CronJobAdapter) Run() {
	_ = c.Run
}

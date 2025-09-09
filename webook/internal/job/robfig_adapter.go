package job

import (
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type RankingJobAdapter struct {
	j       Job
	l       logger.Loggerv1
	summary prometheus.Summary
}

func NewRankingJobAdapter(j Job, l logger.Loggerv1) *RankingJobAdapter {
	return &RankingJobAdapter{
		j: j,
		l: l,
		summary: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: "gobook",
			Subsystem: "webook",
			Name:      "ranking_job_adapter",
			Help:      "ranking job 执行时间统计",
			ConstLabels: map[string]string{
				"name": j.Name(),
			},
		}),
	}
}
func (r *RankingJobAdapter) Run() {
	err := r.j.Run()
	if err != nil {
		r.l.Error("运行任务失败",
			logger.String("job", r.j.Name()),
			logger.Error(err))
	}
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		r.summary.Observe(float64(duration))
	}()
}

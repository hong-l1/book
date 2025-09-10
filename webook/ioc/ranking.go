package ioc

import (
	"github.com/hong-l1/project/webook/internal/job"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/service"
	"github.com/robfig/cron/v3"
)

func InitRankingJob(svc service.RankingService) *job.RankingJob {
	return job.NewRankingJob(svc)
}
func InitJobs(l logger.Loggerv1, rankingjob *job.RankingJob) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	cbd := job.NewCronJobBuilder(l)
	_, err := res.AddJob("0 */3 * * * * ", cbd.Build(rankingjob))
	if err != nil {
		panic(err)
	}
	return res
}

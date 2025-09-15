package ioc

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/job"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/service"
	"time"
)

func InitScheduler(svc service.CronJobService, l logger.Loggerv1, j *job.LocalFuncExecuter) *job.Scheduler {
	res := job.NewScheduler(svc, l)
	res.RegisterExecutor(j)
	return res
}
func InitLocalExecuter(svc service.RankingService) *job.LocalFuncExecuter {
	exe := job.NewLocalFuncExecuter()
	exe.RegisterFunc("ranking_job", func(ctx context.Context, j domain.Job) error {
		ctx, canel := context.WithTimeout(context.Background(), time.Second*30)
		defer canel()
		return svc.TOPN(ctx)
	})
	return exe
}

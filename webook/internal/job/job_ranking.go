package job

import (
	"context"
	"github.com/hong-l1/project/webook/internal/service"
	"time"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, canel := context.WithTimeout(context.Background(), r.timeout)
	defer canel()
	return r.svc.TOPN(ctx)
}

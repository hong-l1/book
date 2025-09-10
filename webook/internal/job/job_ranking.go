package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"github.com/hong-l1/project/webook/internal/service"
	"time"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
	client  *rlock.Client
	key     string
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {

	ctx, canel := context.WithTimeout(context.Background(), r.timeout)
	defer canel()
	r.client.Lock(ctx, r.key, r.timeout)
	return r.svc.TOPN(ctx)
}
func NewRankingJob(svc service.RankingService, client *rlock.Client) *RankingJob {
	return &RankingJob{
		svc:     svc,
		timeout: time.Minute * 30,
		client:  client,
		key:     "rlock_cron_job:ranking",
	}
}

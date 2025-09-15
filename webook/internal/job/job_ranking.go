package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/service"
	"sync"
	"time"
)

type RankingJob struct {
	svc       service.RankingService
	timeout   time.Duration
	client    *rlock.Client
	key       string
	l         logger.Loggerv1
	rlock     *rlock.Lock
	localLock *sync.Mutex
}

func NewRankingJob(svc service.RankingService, client *rlock.Client, l logger.Loggerv1) *RankingJob {
	return &RankingJob{
		svc:       svc,
		timeout:   time.Minute * 30,
		client:    client,
		key:       "rlock_cron_job:ranking",
		l:         l,
		localLock: &sync.Mutex{},
	}
}
func (r *RankingJob) Name() string {
	return "ranking"
}
func (r *RankingJob) Run() error {
	r.localLock.Lock()
	defer r.localLock.Unlock()
	if r.rlock == nil {
		ctx, canel := context.WithTimeout(context.Background(), time.Second)
		defer canel()
		lock, err := r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      0,
		}, time.Second)
		if err != nil {
			return nil
		}
		r.rlock = lock
		go func() {
			err1 := lock.AutoRefresh(r.timeout/2, time.Second)
			if err1 != nil {
				r.l.Error("续约失败", logger.Error(err1))
			}
			r.localLock.Lock()
			r.rlock = nil
			r.localLock.Unlock()
		}()
	}
	ctx, canel := context.WithTimeout(context.Background(), r.timeout)
	defer canel()
	return r.svc.TOPN(ctx)
}

func (r *RankingJob) Close() error {
	r.localLock.Lock()
	lock := r.rlock
	r.rlock = nil
	r.localLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}

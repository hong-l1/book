package service

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/repository"
	"time"
)

type JobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Refresh(ctx context.Context, id int64) error
	ReSetNextTime(ctx context.Context, id int64) error
}
type CronJobService struct {
	repo            repository.JobRepository
	refreshInterval time.Duration
	l               logger.Loggerv1
}

func NewCronJobService(repo repository.JobRepository, refreshInterval time.Duration, l logger.Loggerv1) *CronJobService {
	return &CronJobService{
		repo:            repo,
		refreshInterval: refreshInterval,
		l:               l,
	}
}
func (c *CronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := c.repo.Preempt(ctx)
	ticker := time.NewTicker(c.refreshInterval)
	go func() {
		for range ticker.C {
			ctxfresh, canel := context.WithTimeout(context.Background(), time.Second*1)
			defer canel()
			err1 := c.Refresh(ctxfresh, j.Id)
			if err1 != nil {

			}
		}
	}()
	j.CanelFunc = func() error {
		ticker.Stop()
		ctx, canel := context.WithTimeout(context.Background(), 1*time.Second)
		defer canel()
		return c.repo.Release(ctx, j.Id)
	}
	return j, err
}
func (c *CronJobService) Refresh(ctx context.Context, id int64) error {
	err := c.repo.UpdateUtime(ctx, id)
	if err != nil {
		c.l.Error("续约失败", logger.Int64("jid", id))
	}
	return nil
}

func (c *CronJobService) ReSetNextTime(ctx context.Context, j domain.Job) error {
	next := j.NextTime()
	if next.IsZero() {
		return c.repo.Stop(ctx, j.Id)
	}
	return c.repo.UPdateNextTime(ctx, next, j.Id)
}

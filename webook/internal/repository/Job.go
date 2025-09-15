package repository

import (
	"context"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/repository/dao"
	"time"
)

type JobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UPdateNextTime(ctx context.Context, next time.Time, id int64) error
	Stop(ctx context.Context, id int64) error
}
type CronJobRRepository struct {
	dao dao.CronJobDAO
}

func (r *CronJobRRepository) Stop(ctx context.Context, id int64) {
	//TODO implement me
	panic("implement me")
}

func (r *CronJobRRepository) UPdateNextTime(ctx context.Context, next time.Time, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *CronJobRRepository) UpdateUtime(ctx context.Context, id int64) error {
	return r.dao.UpdateUtime(ctx, id)
}

func (r *CronJobRRepository) Release(ctx context.Context, id int64) error {
	return r.dao.Release(ctx, id)
}

func NewCronJobRRepository(dao dao.CronJobDAO) JobRepository {
	return &CronJobRRepository{dao: dao}
}
func (r *CronJobRRepository) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := r.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		Cfg:  j.Cfg,
		Id:   j.Id,
		Name: j.Name,
	}, nil
}

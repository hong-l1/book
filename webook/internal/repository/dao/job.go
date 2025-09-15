package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Job struct {
	Id       int64 `gorm:"primary_key,auto_increment"`
	Name     string
	Cfg      string
	NextTime int64
	Version  int
	Ctime    int64
	Utime    int64
	Status   int
}
type CronJobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
}
type JobDAO struct {
	db *gorm.DB
}

func (dao *JobDAO) UpdateUtime(ctx context.Context, id int64) error {

}

func (dao *JobDAO) Release(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id=?", id).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  time.Now().UnixMilli(),
	}).Error
}
func NewJobDAO(db *gorm.DB) CronJobDAO {
	return &JobDAO{db: db}
}
func (dao *JobDAO) Preempt(ctx context.Context) (Job, error) {
	db := dao.db.WithContext(ctx)
	for {
		now := time.Now().UnixMilli()
		var job Job
		err := db.WithContext(ctx).Where("status =? AND next_time <=? ", jobStatusWaiting, now).First(&job).Error
		if err != nil {
			return Job{}, err
		}
		res := db.WithContext(ctx).Where("id =? AND version =?", job.Id, job.Version).Model(&Job{}).Updates(map[string]any{
			"status":  jobStatusRunning,
			"utime":   now,
			"version": job.Version + 1,
		})
		if res.Error != nil {
			return Job{}, res.Error
		}
		if res.RowsAffected == 0 {
			continue
		}
		return job, nil
	}
}

const (
	jobStatusWaiting = iota
	jobStatusRunning
	jobStatusPaused
)

package job

import (
	"context"
	"fmt"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/service"
	"time"
)

type Executer interface {
	Name() string
	Exec(ctx context.Context, j domain.Job) error
}
type LocalFuncExecuter struct {
	fn map[string]func(ctx context.Context, j domain.Job) error
}

func NewLocalFuncExecuter() *LocalFuncExecuter {
	return &LocalFuncExecuter{fn: make(map[string]func(ctx context.Context, j domain.Job) error)}
}

func (l *LocalFuncExecuter) RegisterFunc(name string, fn func(ctx context.Context, j domain.Job) error) {
	l.fn[name] = fn
}
func (l *LocalFuncExecuter) Name() string {
	return "LocalFunc"
}

func (l *LocalFuncExecuter) Exec(ctx context.Context, j domain.Job) error {
	fn, ok := l.fn[j.Name]
	if !ok {
		return fmt.Errorf("是否注册了 %s", j.Name)
	}
	return fn(ctx, j)
}

type Scheduler struct {
	svc  service.CronJobService
	l    logger.Loggerv1
	Exes map[string]Executer
}

func NewScheduler(svc service.CronJobService, l logger.Loggerv1) *Scheduler {
	return &Scheduler{svc: svc, l: l}
}

func (s *Scheduler) RegisterExecutor(exec Executer) {
	s.Exes[exec.Name()] = exec
}
func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		dbctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		j, err := s.svc.Preempt(dbctx)
		cancel()
		if err != nil {
			s.l.Error("抢占任务失败", logger.Error(err))
		}
		exe, ok := s.Exes[j.Executer]
		if !ok {
			s.l.Error("没找到执行器", logger.String("executer", j.Executer))
		}
		go func() {
			defer func() {
				err1 := j.CanelFunc()
				s.l.Error("抢占任务失败", logger.Error(err1), logger.Int64("jobid", j.Id))
			}()
			err1 := exe.Exec(ctx, j)
			if err1 != nil {
				s.l.Error("任务执行失败", logger.Error(err1), logger.Int64("jobid", j.Id))
			}
			err1 = s.svc.ReSetNextTime(ctx, j)
			if err1 != nil {

			}
		}()
	}
}

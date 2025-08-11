package failover

import (
	"context"
	"errors"
	"github.com/hong-l1/project/webook/internal/service/sms"
	"log"
	"sync/atomic"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func (f *FailoverSMSService) SendSMS(ctx context.Context, template string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.SendSMS(ctx, template, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("全部都失败了")
}
func (f *FailoverSMSService) SendSMSv1(ctx context.Context, template string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	lenth := uint64(len(f.svcs))
	for i := idx; i < idx+lenth; i++ {
		svc := f.svcs[i]
		err := svc.SendSMS(ctx, template, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			return err
		default:
			log.Println(err)
		}
	}
	return errors.New("全部都失败了")
}
func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

//测试

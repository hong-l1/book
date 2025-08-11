package failover

import (
	"context"
	"github.com/hong-l1/project/webook/internal/service/sms"
	"sync/atomic"
)

type TimeOutFailoverService struct {
	cnt       int32 //当前服务超时数量
	idx       int32 //当前服务节点
	svcs      []sms.Service
	threshold int32 //超时阈值
}

func (t *TimeOutFailoverService) SendSMS(ctx context.Context, template string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt >= t.threshold {
		newidx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newidx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		atomic.AddInt32(&t.idx, 1)
	}
	svc := t.svcs[idx]
	err := svc.SendSMS(ctx, template, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	case nil:
		atomic.AddInt32(&t.cnt, 0)
		return nil
	default:

	}
	return err
}

func NewTimeOutFailoverService(svcs []sms.Service, threshold int32) *TimeOutFailoverService {
	return &TimeOutFailoverService{
		svcs:      svcs,
		threshold: threshold,
	}
}

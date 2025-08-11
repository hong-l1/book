package ratelimit

import (
	"context"
	"errors"
	"github.com/hong-l1/project/webook/internal/pkg/ratelimit"
	"github.com/hong-l1/project/webook/internal/service/sms"
)

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limit
}

func NewService(svc sms.Service, limiter ratelimit.Limit) *RatelimitSMSService {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}
func (s *RatelimitSMSService) SendSMS(ctx context.Context, template string, args []string, numbers ...string) error {
	lim, err := s.limiter.Limited(ctx, "sms:tencent")
	if err != nil {
		return errors.New("限流器异常")
	}
	if lim {
		return errors.New("触发限流")
	}
	err = s.svc.SendSMS(ctx, template, args, numbers...)
	return err
}

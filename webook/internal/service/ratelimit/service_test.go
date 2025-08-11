package ratelimit

import (
	"context"
	"errors"
	"github.com/hong-l1/project/webook/internal/pkg/ratelimit"
	limitmocks "github.com/hong-l1/project/webook/internal/pkg/ratelimit/mocks"
	"github.com/hong-l1/project/webook/internal/service/sms"
	smsmocks "github.com/hong-l1/project/webook/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRatelimitSMSService_SendSMS(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(t *gomock.Controller) (ratelimit.Limit, sms.Service)
		wantErr error
	}{
		{
			name: "发送正常",
			mock: func(ctrl *gomock.Controller) (ratelimit.Limit, sms.Service) {
				limit := limitmocks.NewMockLimit(ctrl)
				sms := smsmocks.NewMockService(ctrl)
				limit.EXPECT().Limited(gomock.Any(), "sms:tencent").Return(false, nil)
				sms.EXPECT().SendSMS(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return limit, sms
			},
			wantErr: nil,
		},
		{
			name: "redis异常",
			mock: func(ctrl *gomock.Controller) (ratelimit.Limit, sms.Service) {
				limit := limitmocks.NewMockLimit(ctrl)
				sms := smsmocks.NewMockService(ctrl)
				limit.EXPECT().Limited(gomock.Any(), "sms:tencent").Return(false, errors.New("1"))
				return limit, sms
			},
			wantErr: errors.New("限流器异常"),
		},
		{
			name: "限流",
			mock: func(ctrl *gomock.Controller) (ratelimit.Limit, sms.Service) {
				limit := limitmocks.NewMockLimit(ctrl)
				sms := smsmocks.NewMockService(ctrl)
				limit.EXPECT().Limited(gomock.Any(), "sms:tencent").Return(true, nil)
				return limit, sms
			},
			wantErr: errors.New("触发限流"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			lim, sms := tc.mock(ctrl)
			s := NewService(sms, lim)
			err := s.SendSMS(context.Background(), "abc", []string{"123"}, "10086")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.uber.org/zap"
)

type Service struct {
	appid     *string
	signature *string
	client    *sms.Client
}

func NewService(appid string, signature string, client *sms.Client) *Service {
	return &Service{
		appid:     ekit.ToPtr[string](appid),
		signature: ekit.ToPtr[string](signature),
		client:    client,
	}
}
func (s *Service) SendSMS(ctx context.Context, biz string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appid
	req.SignName = s.signature
	req.TemplateId = ekit.ToPtr[string](biz)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	resp, err := s.client.SendSms(req)
	zap.L().Debug("发送短信", zap.Any("req", req), zap.Any("resp", resp), zap.Error(err))
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "ok" {
			return fmt.Errorf("发送短信失败%S %S", status.Code, status.Message)
		}
	}
	return nil
}
func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}

package sms

import "context"

// appid,signanature可以一开始初始化好
type Service interface {
	SendSMS(ctx context.Context, biz string, args []string, numbers ...string) error
}

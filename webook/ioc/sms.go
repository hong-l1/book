package ioc

import (
	"github.com/hong-l1/project/webook/internal/service/memory"
	"github.com/hong-l1/project/webook/internal/service/sms"
)

func InitSmsService() sms.Service {
	return memory.NewService()
}

package ioc

import (
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"go.uber.org/zap"
)

func InitLogger() logger.Loggerv1 {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}

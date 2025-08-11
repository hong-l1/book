package service

import (
	"context"
	"fmt"
	"github.com/hong-l1/project/webook/internal/repository"
	"github.com/hong-l1/project/webook/internal/service/sms"
	"math/rand"
)

var ErrSendTooMany = repository.ErrSendTooMany
var ErrVerifyTooMany = repository.ErrVerifyTooMany

const (
	template = "1231241"
)

type CodeService interface {
	Send(ctx context.Context, biz, number string) error
	Verify(ctx context.Context, biz, number, code string) error
}

type codeService struct {
	CodeRepository repository.CodeRepository
	smsSvc         sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{CodeRepository: repo, smsSvc: smsSvc}

}

func (c *codeService) Send(ctx context.Context, biz, number string) error {
	//生成验证码
	//放入redis
	//发送验证码
	code := c.generateCode()
	err := c.CodeRepository.Store(ctx, biz, number, code)
	if err != nil {
		return err
	}
	err = c.smsSvc.SendSMS(ctx, template, []string{code}, number)
	return err
}
func (c *codeService) Verify(ctx context.Context, biz, number, code string) error {
	return c.CodeRepository.Verify(ctx, biz, code, number)
}
func (c *codeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}

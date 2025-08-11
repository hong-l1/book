package service

import (
	"context"
	"fmt"
	"github.com/hong-l1/project/homework/memcache/repository"
	"math/rand"
)

const (
	template = "1231241"
)

type CodeService struct {
	CodeRepository *repository.CodeRepository
	smsSvc         Service
}

func (c *CodeService) Send(ctx context.Context, biz, number string) error {
	//生成验证码
	//放入memcache
	//发送验证码
	code := c.generateCode()
	err := c.CodeRepository.Set(biz, code, number)
	if err != nil {
		return err
	}
	err = c.smsSvc.SendSMS(template, []string{code}, number)
	return err
}
func (c *CodeService) Verify(biz, number, code string) error {
	return c.CodeRepository.Verify(biz, code, number)
}
func (c *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}

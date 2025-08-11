package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hong-l1/project/webook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
	key string
}

func (S *SMSService) SendSMS(ctx context.Context, biz string, args []string, numbers ...string) error {
	var tc Claims
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return S.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return S.svc.SendSMS(ctx, biz, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	TplId string
}

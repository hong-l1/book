package service

import (
	"context"
	"errors"
	"github.com/hong-l1/project/webook/internal/domain"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/repository"
	repomocks "github.com/hong-l1/project/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
)

func Test_userService_LogIn(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		email    string
		password string
		wantUSer domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				r := repomocks.NewMockUserRepository(ctrl)
				r.EXPECT().FindByEmail(context.Background(), domain.User{
					Email:    "3170736324@qq.com",
					Password: "1q2w3e4r",
				}).Return(domain.User{
					Email:    "3170736324@qq.com",
					Password: "$2a$10$qXwzT44xf1Jo8AFKb5fwUuBrW03ffH45rFdJxOKBRUor3RWq69I5C",
				}, nil)
				return r
			},
			email:    "3170736324@qq.com",
			password: "1q2w3e4r",
			wantUSer: domain.User{
				Email:    "3170736324@qq.com",
				Password: "$2a$10$qXwzT44xf1Jo8AFKb5fwUuBrW03ffH45rFdJxOKBRUor3RWq69I5C",
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				r := repomocks.NewMockUserRepository(ctrl)
				r.EXPECT().FindByEmail(context.Background(), domain.User{
					Email:    "3170736324@qq.com",
					Password: "1q2w3e4r",
				}).Return(domain.User{}, repository.ErrUserNotfound)
				return r
			},
			email:    "3170736324@qq.com",
			password: "1q2w3e4r",
			wantUSer: domain.User{},
			wantErr:  repository.ErrUserNotfound,
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				r := repomocks.NewMockUserRepository(ctrl)
				r.EXPECT().FindByEmail(context.Background(), domain.User{
					Email:    "3170736324@qq.com",
					Password: "q2w3e4r",
				}).Return(domain.User{
					Email:    "3170736324@qq.com",
					Password: "$2a$10$qXwzT44xf1Jo8AFKb5fwUuBrW03ffH45rFdJxOKBRUor3RWq69I5C1",
				}, nil)
				return r
			},
			email:    "3170736324@qq.com",
			password: "q2w3e4r",
			wantUSer: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				r := repomocks.NewMockUserRepository(ctrl)
				r.EXPECT().FindByEmail(context.Background(), domain.User{
					Email:    "3170736324@qq.com",
					Password: "1q2w3e4r",
				}).Return(domain.User{}, errors.New("DB错误"))
				return r
			},
			email:    "3170736324@qq.com",
			password: "1q2w3e4r",
			wantUSer: domain.User{},
			wantErr:  errors.New("DB错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl), logger2.NewNopLogger())
			d, err := svc.LogIn(context.Background(), domain.User{
				Email:    tc.email,
				Password: tc.password,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUSer, d)

		})
	}
}
func TestBcrypt(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("1q2w3e4r"), bcrypt.DefaultCost)
	if err == nil {
		log.Println(string(hash))
	}
}

package cache

import (
	"context"
	"errors"
	"github.com/hong-l1/project/webook/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "验证码存储成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), LuaSetCode, []string{"login_" + "10086"}, "123456").Return(res)
				return cmd
			},
			biz:     "login",
			phone:   "10086",
			code:    "123456",
			wantErr: nil,
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("mock error"))
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), LuaSetCode, []string{"login_" + "10086"}, "123456").Return(res)
				return cmd
			},
			biz:     "login",
			phone:   "10086",
			code:    "123456",
			wantErr: errors.New("mock error"),
		},
		{
			name: "验证码发送频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), LuaSetCode, []string{"login_" + "10086"}, "123456").Return(res)
				return cmd
			},
			biz:     "login",
			phone:   "10086",
			code:    "123456",
			wantErr: ErrSendTooMany,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(100))
				cmd.EXPECT().Eval(gomock.Any(), LuaSetCode, []string{"login_" + "10086"}, "123456").Return(res)
				return cmd
			},
			biz:     "login",
			phone:   "10086",
			code:    "123456",
			wantErr: errors.New("system error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCodeCache(tc.mock(ctrl))
			err := c.Set(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

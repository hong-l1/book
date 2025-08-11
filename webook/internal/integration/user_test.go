package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hong-l1/project/webook/internal/integration/startup"
	"github.com/hong-l1/project/webook/internal/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandle_SendLoginSMSCode(t *testing.T) {
	server := startup.InitWebServer()
	rdb := startup.InitRedis()
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		reqBody  web.Result
		wantCode int
		phone    string
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: web.Result{
				Msg: "发送成功",
			},
			phone:    "15223456789",
			wantCode: 200,
		},
		{
			name: "发送频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*50).Err()
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			reqBody: web.Result{
				Msg: "发送太频繁，请稍后再试",
			},
			phone: "15212345678",
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				err := rdb.Set(ctx, key, "123456", 0).Err()
				require.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				key := "phone_code:login:15212345678"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			reqBody: web.Result{
				Msg:  "系统错误",
				Code: 5,
			},
			phone: "15212345678",
		},
		{
			name: "为输入手机号",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: web.Result{
				Msg:  "请输入合法的手机号",
				Code: 4,
			},
			phone: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewBuffer([]byte(fmt.Sprintf(`{"phone": "%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp := httptest.NewRecorder()
			var webresp web.Result
			server.ServeHTTP(resp, req)
			err = json.NewDecoder(resp.Body).Decode(&webresp)
			assert.Equal(t, tc.reqBody, webresp)
			tc.after(t)
		})
	}
}

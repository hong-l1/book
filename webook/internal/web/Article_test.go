package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hong-l1/project/webook/internal/domain"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/service"
	svcmocks "github.com/hong-l1/project/webook/internal/service/mocks"
	ijwt "github.com/hong-l1/project/webook/internal/web/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandle_Publish(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) service.ArticleService
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantBody   Result
		reqbody    string
	}{
		{
			name: "新建并发表",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				a := svcmocks.NewMockArticleService(ctrl)
				a.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return a
			},
			reqbody: `{
    "title": "我的标题",
    "content": "我的内容"  
}`,
			wantCode: 200,
			wantBody: Result{
				Msg:  "ok",
				Data: float64(1),
			},
		},
		{
			name: "publish失败",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				a := svcmocks.NewMockArticleService(ctrl)
				a.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("publish fail"))
				return a
			},
			reqbody: `{
    "title": "我的标题",
    "content": "我的内容"  
}`,
			wantCode: 200,
			wantBody: Result{
				Msg:  "系统错误",
				Code: 5,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			h := NewArticleHandle(logger2.NewNopLogger(), tc.mock(ctrl))
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("claim", &ijwt.Claim{
					UserId: 123,
				})
			})
			h.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBuffer([]byte(tc.reqbody)))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code == http.StatusOK {
				return
			}
			var webresp Result
			err = json.NewDecoder(resp.Body).Decode(&webresp)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webresp)
		})
	}
}

package repository

import (
	"context"
	"database/sql"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/repository/cache"
	cachemocks "github.com/hong-l1/project/webook/internal/repository/cache/mocks"
	"github.com/hong-l1/project/webook/internal/repository/dao"
	daomocks "github.com/hong-l1/project/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，但是查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().GetUserCache(gomock.Any(), domain.User{Id: int64(123)}).Return(domain.User{}, cache.ErrKeyNotFound)
				d := daomocks.NewMockUserDao(ctrl)
				d.EXPECT().FindById(gomock.Any(), domain.User{Id: int64(123)}).Return(dao.User{
					Email: sql.NullString{
						String: "3170736324@qq.com",
						Valid:  true,
					},
					Id:       int64(123),
					Password: "this is password",
					Phone: sql.NullString{
						String: "18818188188188188188188",
						Valid:  true,
					},
				}, nil)
				c.EXPECT().SetUserCache(gomock.Any(), domain.User{
					Email:    "3170736324@qq.com",
					Id:       int64(123),
					Password: "this is password",
					Phone:    "18818188188188188188188",
				}).Return(nil)
				return d, c
			},
			id: 123,
			wantUser: domain.User{
				Email:    "3170736324@qq.com",
				Id:       int64(123),
				Password: "this is password",
				Phone:    "18818188188188188188188",
			},
			wantErr: nil,
		},
		{
			name: "缓存命中,查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().GetUserCache(gomock.Any(), domain.User{Id: int64(123)}).Return(domain.User{
					Email:    "3170736324@qq.com",
					Id:       int64(123),
					Password: "this is password",
					Phone:    "18818188188188188188188",
				}, nil)
				d := daomocks.NewMockUserDao(ctrl)
				return d, c
			},
			id: 123,
			wantUser: domain.User{
				Email:    "3170736324@qq.com",
				Id:       int64(123),
				Password: "this is password",
				Phone:    "18818188188188188188188",
			},
			wantErr: nil,
		},
		{
			name: "缓存和数据库均无",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().GetUserCache(gomock.Any(), domain.User{Id: int64(123)}).Return(domain.User{}, cache.ErrKeyNotFound)
				d := daomocks.NewMockUserDao(ctrl)
				d.EXPECT().FindById(gomock.Any(), domain.User{Id: int64(123)}).Return(dao.User{}, ErrUserNotfound)
				return d, c
			},
			id:       123,
			wantUser: domain.User{},
			wantErr:  ErrUserNotfound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			d, c := tc.mock(ctrl)
			repo := NewRepository(d, c)
			u, err := repo.FindById(context.Background(), domain.User{
				Id: tc.id,
			})
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

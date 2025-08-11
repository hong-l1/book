package repository

import (
	"context"
	"database/sql"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/repository/cache"
	"github.com/hong-l1/project/webook/internal/repository/dao"
)

var (
	ErrUserNotfound   = dao.ErrUserNotfound
	ErrUserDuplicated = dao.ErrUserDuplicated
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, user domain.User) (domain.User, error)
	Edit(ctx context.Context, user domain.User) error
	Profile(ctx context.Context, user domain.User) (domain.User, error)
	FindById(ctx context.Context, user domain.User) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
	}
}
func (ur *CacheUserRepository) Create(ctx context.Context, user domain.User) error {
	return ur.dao.Insert(ctx, ur.DomainToentity(user))
}
func (ur *CacheUserRepository) FindByEmail(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDomain(u), nil
}
func (ur *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDomain(u), nil
}
func (ur *CacheUserRepository) Edit(ctx context.Context, user domain.User) error {
	err := ur.dao.Edit(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
func (ur *CacheUserRepository) Profile(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := ur.dao.FindById(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Birthday:     u.Birthday,
		Nickname:     u.Nickname,
		Introduction: u.Introduction,
	}, nil
}
func (ur *CacheUserRepository) FindById(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := ur.cache.GetUserCache(ctx, user)
	if err == nil {
		//必然有数据
		return u, nil
	}
	//没有找到数据
	temp, err := ur.dao.FindById(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	u = ur.entityToDomain(temp)
	err = ur.cache.SetUserCache(ctx, u)
	return u, err
	//if errors.Is(err, cache.ErrKeyNotFound) {
	//	//取数据库里面加载
	//}
	//别的错误，reids可能崩了
	//1.取数据库加载，需要保护数据库（限流）
	//2.不加载数据，影响一点用户体验
}
func (ur *CacheUserRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := ur.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDomain(u), nil
}
func (ur *CacheUserRepository) entityToDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Password: user.Password,
		Birthday: user.Birthday,
		Phone:    user.Phone.String,
		WechatInfo: domain.WechatInfo{
			OpenID:  user.WechatOpenId.String,
			UnionId: user.WechatUnionId.String,
		},
	}
}
func (ur *CacheUserRepository) DomainToentity(user domain.User) dao.User {
	return dao.User{
		Id: user.Id,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Password: user.Password,
		Birthday: user.Birthday,
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		WechatOpenId: sql.NullString{
			String: user.WechatInfo.OpenID,
			Valid:  user.WechatInfo.OpenID != "",
		},
		WechatUnionId: sql.NullString{
			String: user.WechatInfo.UnionId,
			Valid:  user.WechatInfo.UnionId != "",
		},
	}
}

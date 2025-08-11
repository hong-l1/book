package service

import (
	"context"
	"errors"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserDuplicated = repository.ErrUserDuplicated
var ErrInvalidUserOrPassword = errors.New("账号或密码不对")

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	LogIn(ctx context.Context, user domain.User) (domain.User, error)
	Edit(ctx context.Context, user domain.User) error
	Profile(ctx context.Context, user domain.User) (domain.User, error)
	FindORCreate(ctx context.Context, phone string) (domain.User, error)
	FindORCreateBywechat(ctx context.Context, wechat domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	l    logger.Loggerv1
}

func NewUserService(repo repository.UserRepository, l logger.Loggerv1) UserService {
	return &userService{
		repo: repo,
		l:    l,
	}
}
func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	//加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}
func (svc *userService) LogIn(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, user)
	if errors.Is(err, repository.ErrUserNotfound) {
		return domain.User{}, gorm.ErrRecordNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
func (svc *userService) Edit(ctx context.Context, user domain.User) error {
	err := svc.repo.Edit(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
func (svc *userService) Profile(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := svc.repo.Profile(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
func (svc *userService) FindORCreateBywechat(ctx context.Context, wechat domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, wechat.OpenID)
	if !errors.Is(err, repository.ErrUserNotfound) {
		return u, err
	}
	u = domain.User{
		WechatInfo: wechat,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil || !errors.Is(err, ErrUserDuplicated) {
		return u, err
	}
	return svc.repo.FindByWechat(ctx, wechat.OpenID)
}
func (svc *userService) FindORCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, repository.ErrUserNotfound) {
		return u, err
	}
	svc.l.Info("用户未注册", logger.String("phone", phone))
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil || !errors.Is(err, ErrUserDuplicated) {
		return u, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}

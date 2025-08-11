package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/hong-l1/project/webook/internal/domain"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserNotfound   = gorm.ErrRecordNotFound
	ErrUserDuplicated = errors.New("冲突")
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, user domain.User) (User, error)
	FindById(ctx context.Context, user domain.User) (User, error)
	Edit(ctx context.Context, user domain.User) error
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
}
type GormUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GormUserDao{db: db}
}
func (ud *GormUserDao) Insert(ctx context.Context, u User) error {
	time := time.Now().UnixMilli()
	u.Utime = time
	u.Ctime = time
	err := ud.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictErrNo {
			return ErrUserDuplicated
		}
	}
	return err
}
func (ud *GormUserDao) FindByEmail(ctx context.Context, user domain.User) (User, error) {
	var temp User
	return temp, ud.db.WithContext(ctx).First(&temp, "email = ?", user.Email).Error
}
func (ud *GormUserDao) Edit(ctx context.Context, user domain.User) error {
	u, _ := ud.FindById(ctx, user)
	return ud.db.Model(&u).Updates(map[string]interface{}{
		"Nickname":     user.Nickname,
		"Birthday":     user.Birthday,
		"Introduction": user.Introduction,
	}).Error
}
func (ud *GormUserDao) FindById(ctx context.Context, user domain.User) (User, error) {
	var temp User
	return temp, ud.db.WithContext(ctx).First(&temp, "id = ?", user.Id).Error
}
func (ud *GormUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var temp User
	return temp, ud.db.WithContext(ctx).First(&temp, "phone = ?", phone).Error
}
func (ud *GormUserDao) FindByWechat(ctx context.Context, openId string) (User, error) {
	var temp User
	return temp, ud.db.WithContext(ctx).First(&temp, "Wechat_Open_Id= ?", openId).Error
}

type User struct {
	Id            int64          `gorm:"primary_key,auto_increment"`
	Email         sql.NullString `gorm:"unique"`
	Password      string
	Ctime         int64
	Utime         int64
	Nickname      string
	Birthday      string
	Introduction  string
	Phone         sql.NullString `gorm:"unique"`
	WechatUnionId sql.NullString
	WechatOpenId  sql.NullString `gorm:"unique"`
}

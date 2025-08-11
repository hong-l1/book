package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		user    User
		wantErr error
		mock    func(t *testing.T) *sql.DB
	}{
		{
			name: "插入成功",
			user: User{
				Email: sql.NullString{
					String: "3170736324@qq.com",
				},
			},
			wantErr: nil,
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(123, 1)
				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectExec("INSERT INTO .*").
					WillReturnResult(mockRes)
				return db
			},
		},
		{
			name: "邮箱冲突",
			user: User{
				Email: sql.NullString{
					String: "3170736324@qq.com",
				},
			},
			wantErr: ErrUserDuplicated,
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(&mysql.MySQLError{
						Number: 1062,
					})
				return db
			},
		},
		{
			name: "数据库错误",
			user: User{
				Email: sql.NullString{
					String: "3170736324@qq.com",
				},
			},
			wantErr: errors.New("数据库错误"),
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(errors.New("数据库错误"))
				return db
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(gorm_mysql.New(gorm_mysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			require.NoError(t, err)
			d := NewUserDao(db)
			err = d.Insert(context.Background(), tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

package dao

import (
	"github.com/hong-l1/project/webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &article.Article{})
}

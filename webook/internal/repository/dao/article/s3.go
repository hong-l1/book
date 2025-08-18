package article

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

var accessKey = "BIWKkyBa1iJ-NEVPx0rSj4ol0fJBrF_F8rSPoiL8"
var secretKey = "OCUZEvmfyZ-_aThggjeAtCau_g2Mfzq8gPUnHYN7"
var bucket = "mybookdemo"
var domain = "http://t135by11b.hn-bkt.clouddn.com"

type S3DAO struct {
	oss *auth.Credentials
	GORMArticleDao
	cfg storage.Config
}

func NewOssDAO(db *gorm.DB) ArticleDao {
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan, // 你可以外部传参决定区域
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	return &S3DAO{
		oss: mac,
		GORMArticleDao: GORMArticleDao{
			db: db,
		},
		cfg: cfg,
	}
}
func (o *S3DAO) Sync(ctx context.Context, art Article) (int64, error) {
	var id = art.Id
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		now := time.Now().UnixMilli()
		if id == 0 {
			id, err = o.GORMArticleDao.Insert(ctx, art)
		} else {
			err = o.GORMArticleDao.UpdateById(ctx, art)
		}
		if err != nil {
			return err
		}
		art.Id = id
		publishArt := PublishArticleDAO{
			Article{
				Id:       art.Id,
				AuthorId: art.AuthorId,
				Title:    art.Title,
				Ctime:    now,
				Utime:    now,
			},
		}
		publishArt.Content = ""
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":  publishArt.Title,
				"utime":  publishArt.Utime,
				"status": publishArt.Status,
			}),
		}).Create(&publishArt).Error
	})
	if err != nil {
		return id, err
	}
	contenturl, err := o.Upload(art.Content)
	if err != nil {
		return id, err
	}
	err = o.db.WithContext(ctx).Model(&PublishArticleDAO{}).
		Where("id = ?", id).
		Update("content", contenturl).Error
	return id, err
}
func (o *S3DAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	panic("implement me")
}
func (s *S3DAO) Upload(content string) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	uploadToken := putPolicy.UploadToken(s.oss)
	formUploader := storage.NewFormUploader(&s.cfg)
	ret := storage.PutRet{}
	objectKey := fmt.Sprintf("articles/%d.txt", time.Now().UnixNano())
	putExtra := &storage.PutExtra{
		MimeType: "text/plain; charset=utf-8",
	}
	err := formUploader.Put(context.Background(), &ret, uploadToken, objectKey,
		strings.NewReader(content), int64(len(content)), putExtra)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", domain, objectKey), nil
}

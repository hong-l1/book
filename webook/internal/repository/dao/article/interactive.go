package article

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type Interactive struct {
	ID         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_id_type"`
	Biz        string `gorm:"uniqueIndex:biz_id_type;type:varchar(128)"`
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Ctime      int64
	Utime      int64
}
type UserLike struct {
	ID    int64  `gorm:"primaryKey,autoIncrement"`
	Biz   string `gorm:"uniqueIndex:uid_biz_id_type;type:varchar(128)"`
	BizId int64  `gorm:"uniqueIndex:uid_biz_id_type"`
	Uid   int64  `gorm:"uniqueIndex:uid_biz_id_type"`
	Ctime int64
	Utime int64
	//软删除
	Status uint8
}
type UserCollectionBiz struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Uid   int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`
	Cid   int64  `gorm:"index"`
	Utime int64
	Ctime int64
}
type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizid int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, bizId int64, uid int64) error
	InsertCollectionBiz(ctx context.Context, ob UserCollectionBiz) error
	GetLikeInfo(ctx context.Context, biz string, artid int64, uid int64) (interface{}, interface{})
	GetCollectInfo(ctx context.Context, biz string, artid int64, uid int64) (interface{}, interface{})
	Get(ctx context.Context, biz string, artid int64) (Interactive, error)
}
type GORMInteractiveDAO struct {
	db *gorm.DB
}

func (g *GORMInteractiveDAO) GetLikeInfo(ctx context.Context, biz string, artid int64, uid int64) (interface{}, interface{}) {
	//TODO implement me
	panic("implement me")
}

func (g *GORMInteractiveDAO) GetCollectInfo(ctx context.Context, biz string, artid int64, uid int64) (interface{}, interface{}) {
	//TODO implement me
	panic("implement me")
}

func (g *GORMInteractiveDAO) Get(ctx context.Context, biz string, artid int64) (Interactive, error) {
	var result Interactive
	err := g.db.WithContext(ctx).Where("biz=? and article_id=?", biz, artid).First(&result).Error
	return result, err
}

func (g *GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context, ob UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	ob.Ctime = now
	ob.Utime = now
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&ob).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"utime":       now,
				"Collect_cnt": gorm.Expr("Collect_Cnt + ?", 1),
			}),
		}).Create(&Interactive{
			Biz:        ob.Biz,
			BizId:      ob.BizId,
			Utime:      now,
			Ctime:      now,
			CollectCnt: 1,
		}).Error
	})
}

func (g *GORMInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId int64, uid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLike{}).Where("biz = ? AND uid = ? AND biz_id = ?", biz, uid, bizId).Updates(map[string]interface{}{
			"utime":  now,
			"status": 0,
		}).Error
		if err != nil {
			return err
		}
		return tx.Model(&Interactive{}).Where("biz_id = ? AND biz", bizId, biz).Updates(map[string]interface{}{
			"utime":    now,
			"like_cnt": gorm.Expr("like_cnt - ?", 1),
		}).Error
	})
}

func (g *GORMInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, bizid int64, uid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"utime":  now,
				"status": 1,
			}),
		}).Create(&UserLike{
			Biz:    biz,
			BizId:  bizid,
			Uid:    uid,
			Status: 1,
			Ctime:  now,
			Utime:  now,
		}).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt": gorm.Expr("like_cnt + ?", 1),
				"utime":    time.Now().UnixMilli(),
			}),
		}).Create(&Interactive{
			BizId:   bizid,
			Biz:     biz,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
}
func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{db: db}
}
func (g *GORMInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_cnt": gorm.Expr("read_cnt + ?", 1),
			"utime":    time.Now().UnixMilli(),
		}),
	}).Create(&Interactive{
		BizId:   bizId,
		Biz:     biz,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

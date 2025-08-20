package article

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBDAO struct {
	Client   *mongo.Client
	Database *mongo.Database
	//制作库
	Docol *mongo.Collection
	//线上库
	Linecol *mongo.Collection
	Node    *snowflake.Node
}

func (m *MongoDBDAO) GetbyAuthor(ctx context.Context, id int64, offset int, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) GetById(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	id := m.Node.Generate().Int64()
	_, err := m.Docol.InsertOne(ctx, art)
	return id, err
}
func (m *MongoDBDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	filter := bson.M{
		"id":        art.Id,
		"author_id": art.AuthorId,
	}
	update := bson.M{"$set": bson.M{
		"utime":   now,
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
	}}
	res, err := m.Docol.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("文档不存在（id: %d, author_id: %d）", art.Id, art.AuthorId)
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("数据未变更（可能新值和旧值相同）")
	}
	return nil
}
func (m *MongoDBDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	now := time.Now().UnixMilli()
	filter := bson.M{"id": id}
	update := bson.M{"$set": PublishArticleDAO{art}, "$setOnInsert": bson.M{"Ctime": now}}
	_, err = m.Linecol.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return id, err
}
func (m *MongoDBDAO) Upsert(ctx context.Context, art PublishArticleDAO) error {
	//TODO implement me
	panic("implement me")
}
func (m *MongoDBDAO) SyncStatus(ctx context.Context, articleId int64, AuthorId int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}
func InitCollections(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "ctime", Value: 1},
			},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("article").Indexes().
		CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("published_article").Indexes().
		CreateMany(ctx, index)
	return err
}
func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) ArticleDao {
	return &MongoDBDAO{
		Database: db,
		Docol:    db.Collection("article"),
		Linecol:  db.Collection("publish_article"),
		Node:     node,
	}
}

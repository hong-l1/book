package mongo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}

func TestMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		// 每个命令（查询）执行之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		// 执行成功
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
		},
		// 执行失败
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
		},
	}
	opts := options.Client().
		ApplyURI("mongodb://localhost:27017/book").
		SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)
	mdb := client.Database("webook")
	col := mdb.Collection("articles")
	defer func() {
		_, err = col.DeleteMany(ctx, bson.M{})
	}()
	res, err := col.InsertOne(ctx, Article{
		Id:      123,
		Title:   "我的标题",
		Content: "我的内容",
		Status:  0,
	})
	assert.NoError(t, err)
	fmt.Printf("Inserted ID: %d \n", res.InsertedID)
	filter := bson.M{"id": 123}
	var result Article
	err = col.FindOne(ctx, filter).Decode(&result)
	assert.NoError(t, err)
	fmt.Println("找到用户:", result)
	update := bson.M{
		"$set": bson.M{
			"status": 1,
		},
	}
	updateres, err := col.UpdateOne(ctx, filter, update)
	assert.NoError(t, err)
	fmt.Println(updateres.MatchedCount, updateres.ModifiedCount)
	filterAnd := bson.M{
		"$or": bson.A{
			bson.M{"title": "Golang"},
			bson.M{"title": "MongoDB"},
		},
	}
	filterOr := bson.M{
		"$and": bson.A{
			bson.M{"title": "Golang"},
			bson.M{"title": "MongoDB"},
		},
	}
	filterIn := bson.M{
		"status": bson.M{"$in": bson.A{0, 1}},
	}
	_, _ = col.UpdateOne(ctx, filterAnd, update)
	_, _ = col.UpdateOne(ctx, filterOr, update)
	_, _ = col.UpdateOne(ctx, filterIn, update)
}

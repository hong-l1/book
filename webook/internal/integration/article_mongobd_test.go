package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/integration/startup"
	"github.com/hong-l1/project/webook/internal/repository/dao/article"
	ijwt "github.com/hong-l1/project/webook/internal/web/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type ArticleMongoHandlerTestSuite struct {
	suite.Suite
	server  *gin.Engine
	mdb     *mongo.Database
	col     *mongo.Collection
	liveCol *mongo.Collection
}

func (s *ArticleMongoHandlerTestSuite) SetupSuite() {
	s.server = gin.Default()
	s.server.Use(func(context *gin.Context) {
		context.Set("claim", &ijwt.Claim{
			UserId: 123,
		})
		context.Next()
	})
	s.mdb = startup.InitMongoDB()
	s.col = s.mdb.Collection("articles")
	s.liveCol = s.mdb.Collection("published_articles")
	node, err := snowflake.NewNode(1)
	assert.NoError(s.T(), err)
	hdl := startup.InitArticleHandle(article.NewMongoDBDAO(s.mdb, node))
	hdl.RegisterRoutes(s.server)
}
func TestMongoArticle(t *testing.T) {
	suite.Run(t, new(ArticleMongoHandlerTestSuite))
}
func (s *ArticleMongoHandlerTestSuite) TestCleanMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := s.col.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
	_, err = s.liveCol.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
}
func (s *ArticleMongoHandlerTestSuite) TestArticleHandler_Edit() {
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		art      Article
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "新建帖子",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				var art article.Article
				err := s.col.FindOne(ctx, bson.M{"author_id": 123}).Decode(&art)
				assert.NoError(s.T(), err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.AuthorId > 0)
				art.Ctime = 0
				art.Utime = 0
				art.AuthorId = 0
				temp := article.Article{
					Ctime:    0,
					Utime:    0,
					Title:    "hello，你好",
					Content:  "随便试试",
					AuthorId: 123,
					Status:   2,
				}
				assert.Equal(s.T(), art, temp)
			},
			art: Article{
				Title:   "hello，你好",
				Content: "随便试试",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 1,
			},
		},
		{
			name: "更新帖子",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := s.col.InsertOne(ctx, bson.M{
					"author_id": 123,
					"title":     "我的标题",
					"content":   "我的内容",
					"id":        2,
					"status":    domain.ArticleStatusUnpublished.ToUint8(),
				})
				assert.NoError(s.T(), err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				var art article.Article
				err := s.col.FindOne(ctx, bson.M{"author_id": 123}).Decode(&art)
				assert.NoError(s.T(), err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.AuthorId > 0)
				art.Ctime = 0
				art.Utime = 0
				art.AuthorId = 0
				temp := article.Article{
					Ctime:    0,
					Utime:    0,
					Title:    "hello，你好",
					Content:  "随便试试",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}
				assert.Equal(s.T(), art, temp)
			},
			art: Article{
				Title:   "hello，你好",
				Content: "随便试试",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 1,
			},
		},
		{
			name: "更新别人的帖子",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err := s.col.InsertOne(ctx, bson.M{
					"author_id": 321,
					"title":     "我的标题",
					"content":   "我的内容",
					"id":        2,
					"status":    domain.ArticleStatusUnpublished.ToUint8(),
				})
				assert.NoError(s.T(), err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				var art article.Article
				err := s.col.FindOne(ctx, bson.M{"author_id": 123}).Decode(&art)
				assert.NoError(s.T(), err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				assert.True(t, art.AuthorId > 0)
				art.Ctime = 0
				art.Utime = 0
				art.AuthorId = 0
				temp := article.Article{
					Ctime:    0,
					Utime:    0,
					Title:    "hello，你好",
					Content:  "随便试试",
					AuthorId: 321,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}
				assert.Equal(s.T(), art, temp)
			},
			art: Article{
				Title:   "hello，你好",
				Content: "随便试试",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 1,
			},
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			reqbody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqbody))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code == http.StatusOK {
				return
			}
			var webresp Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webresp)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webresp)
			tc.after(t)
		})
	}
}
func (s *ArticleMongoHandlerTestSuite) TestArticleHandler_Publish() {
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		art      Article
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "新建帖子并发表",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var art article.Article
				err := s.col.FindOne(ctx, bson.M{"author_id": 123}).Decode(&art)
				assert.NoError(t, err)
				assert.True(t, art.Id > 0)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Id = 0
				art.Utime = 0
				art.Ctime = 0
				temp := article.Article{
					Title:    "hello，你好",
					Content:  "随便试试",
					Ctime:    0,
					Utime:    0,
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}
				assert.Equal(t, art, temp)
				var publishedArt article.PublishArticleDAO
				err = s.liveCol.FindOne(ctx, bson.M{"author_id": 123}).Decode(&publishedArt)
				assert.NoError(t, err)
				assert.True(t, publishedArt.Id > 0)
				assert.True(t, publishedArt.Ctime > 0)
				assert.True(t, publishedArt.Utime > 0)
				publishedArt.Id = 0
				publishedArt.Utime = 0
				publishedArt.Ctime = 0
				assert.Equal(t, publishedArt, temp)
			},
			art: Article{
				Title:   "hello，你好",
				Content: "随便试试",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 1,
			},
		},
		{
			name: "修改帖子并发表",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				templ := bson.M{
					"author_id": int64(123),
					"id":        int64(1),
					"title":     "<UNK>",
					"content":   "<UNK>",
					"status":    domain.ArticleStatusUnpublished.ToUint8(),
				}
				_, err := s.col.InsertOne(ctx, templ)
				if err != nil {
					panic(err)
				}
				_, err = s.liveCol.InsertOne(ctx, templ)
				if err != nil {
					panic(err)
				}
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				var art article.Article
				err := s.col.FindOne(ctx, bson.M{"author_id": 123}).Decode(&art)
				assert.NoError(t, err)
				assert.True(t, art.Id > 0)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Id = 0
				art.Utime = 0
				art.Ctime = 0
				temp := article.Article{
					Id:       int64(1),
					Title:    "hello，你好",
					Content:  "随便试试",
					Ctime:    0,
					Utime:    0,
					AuthorId: int64(123),
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}
				assert.Equal(t, art, temp)
				var publishedArt article.PublishArticleDAO
				err = s.liveCol.FindOne(ctx, bson.M{"author_id": 123}).Decode(&publishedArt)
				assert.NoError(t, err)
				assert.True(t, publishedArt.Id > 0)
				assert.True(t, publishedArt.Ctime > 0)
				assert.True(t, publishedArt.Utime > 0)
				publishedArt.Id = 0
				publishedArt.Utime = 0
				publishedArt.Ctime = 0
				assert.Equal(t, publishedArt, temp)
			},
			art: Article{
				Id:      1,
				Title:   "hello，你好",
				Content: "随便试试",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "ok",
				Data: 1,
			},
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			reqbody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewBuffer(reqbody))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code == http.StatusOK {
				return
			}
			var webresp Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webresp)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webresp)
			tc.after(t)
		})
	}
}

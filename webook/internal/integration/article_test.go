package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hong-l1/project/webook/internal/integration/startup"
	"github.com/hong-l1/project/webook/internal/repository/dao/article"
	ijwt "github.com/hong-l1/project/webook/internal/web/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

//测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func (suite *ArticleTestSuite) SetupSuite() {
	suite.server = gin.Default()
	suite.db = startup.InitDB()
	suite.server.Use(func(c *gin.Context) {
		c.Set("claim", &ijwt.Claim{
			UserId: 123,
		})
	})
	artHal := startup.InitArticleHandle()
	artHal.RegisterRoutes(suite.server)
}
func (s *ArticleTestSuite) TearDownTest() {
	s.db.Exec("truncate table `articles`")
}

func (suite *ArticleTestSuite) TestEdit() {
	t := suite.T()
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		art      Article
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "帖子保存成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var art article.Article
				err := suite.db.Where("id= ?", 1).First(&art).Error
				require.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       1,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
				}, art)
			},
			art: Article{
				Title:   "我的标题",
				Content: "我的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: int64(1),
				Msg:  "ok",
			},
		},
		{
			name: "修改已有帖子",
			before: func(t *testing.T) {
				err := suite.db.Create(&article.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Ctime:    123,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := suite.db.Where("id= ?", 2).First(&art).Error
				require.NoError(t, err)
				assert.True(t, art.Utime > 234)
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					Ctime:    123,
				}, art)
			},
			art: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: int64(2),
				Msg:  "ok",
			},
		},
		{
			name: "修改别人的帖子",
			before: func(t *testing.T) {
				err := suite.db.Create(&article.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 789,
					Ctime:    123,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art article.Article
				err := suite.db.Where("id= ?", 3).First(&art).Error
				require.NoError(t, err)
				assert.Equal(t, article.Article{
					Id:       3,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 789,
					Ctime:    123,
					Utime:    234,
				}, art)
			},
			art: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Data: 0,
				Msg:  "系统错误",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			reqbody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqbody))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			resp := httptest.NewRecorder()
			suite.server.ServeHTTP(resp, req)
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
func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

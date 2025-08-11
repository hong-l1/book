package service

import (
	"context"
	"errors"
	"github.com/hong-l1/project/webook/internal/domain"
	logger2 "github.com/hong-l1/project/webook/internal/pkg/logger"
	"github.com/hong-l1/project/webook/internal/repository/article"
	repomocks2 "github.com/hong-l1/project/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestServiceArticle_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (article.ArticleReaderRepository, article.ArticleAuthorRepository)
		art     domain.Article
		wantErr error
		wantid  int64
	}{
		{
			name: "新建发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleReaderRepository, article.ArticleAuthorRepository) {
				author := repomocks2.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					}}).Return(int64(1), nil)
				reader := repomocks2.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					}}).Return(nil)
				return reader, author
			},
			art: domain.Article{
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: nil,
			wantid:  int64(1),
		},
		{
			name: "修改并发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleReaderRepository, article.ArticleAuthorRepository) {
				author := repomocks2.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					}}).Return(nil)
				reader := repomocks2.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					}}).Return(nil)
				return reader, author
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: nil,
			wantid:  int64(2),
		},
		{
			//新建保存失败
			//创建保存失败
			name: "保存到制作库失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleReaderRepository, article.ArticleAuthorRepository) {
				author := repomocks2.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					}}).Return(errors.New("mock db error"))
				reader := repomocks2.NewMockArticleReaderRepository(ctrl)
				return reader, author
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: errors.New("mock db error"),
			wantid:  int64(0),
		},
		{
			name: "保存制作库成功，重试到线上库成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleReaderRepository, article.ArticleAuthorRepository) {
				author := repomocks2.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					}}).Return(nil)
				reader := repomocks2.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					}}).Return(int64(0), errors.New("mock db error"))
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					}}).Return(int64(2), nil)
				return reader, author
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: nil,
			wantid:  int64(2),
		},
		{
			name: "保存制作库成功，重试全部失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleReaderRepository, article.ArticleAuthorRepository) {
				author := repomocks2.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					}}).Return(nil)
				reader := repomocks2.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					}}).Return(int64(0), errors.New("mock db error")).Times(3)
				return reader, author
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: errors.New("mock db error"),
			wantid:  int64(0),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			reader, author := tc.mock(ctrl)
			svc := NewServiceArticlev1(author, reader, &logger2.NopLogger{})
			id, err := svc.Publishv1(context.Background(), tc.art)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantid, id)
		})
	}
}

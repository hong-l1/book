package web

import "github.com/hong-l1/project/webook/internal/domain"

type Req struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
type LikeReq struct {
	Id   int64
	Like bool `json:"like"`
}
type CollectReq struct {
	id int64
}

func (req Req) todomain(userid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: userid,
		},
	}
}

type ArticleVO struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Abstract   string `json:"abstract"`
	Content    string `json:"content"`
	Status     uint8  `json:"status"`
	AuthorId   int64  `json:"authorId"`
	AuthorName string `json:"authorName"`

	ReadCnt    int64  `json:"read_Cnt"`
	LikeCnt    int64  `json:"like_Cnt"`
	CollectCnt int64  `json:"collect_Cnt"`
	Liked      bool   `json:"liked"`
	Collected  bool   `json:"collected"`
	Ctime      string `json:"ctime"`
	Utime      string `json:"utime"`
}

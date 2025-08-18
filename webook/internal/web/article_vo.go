package web

import "github.com/hong-l1/project/webook/internal/domain"

type Req struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	id      int64
}
type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
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
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Content  string `json:"content"`
	Status   uint8  `json:"status"`
	Author   string `json:"author"`
	Ctime    string `json:"ctime"`
	Utime    string `json:"utime"`
}

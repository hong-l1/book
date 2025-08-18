package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
	Ctime   time.Time
	Utime   time.Time
}

type ArticleStatus uint8

const (
	ArticleStatusUnkown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

type Author struct {
	Id   int64
	Name string
}

func (a Article) Abstract() string {
	cs := []rune(a.Content)
	if len(cs) <= 100 {
		return string(cs)
	}
	return string(cs[:100])
}

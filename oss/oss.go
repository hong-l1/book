package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Article struct {
	Id         string `bson:"_id,omitempty"`
	Title      string `bson:"title"`
	AuthorId   string `bson:"author_id"`
	ContentURL string `bson:"content_url"`
	CreatedAt  int64  `bson:"created_at"`
}

// 上传内容到七牛云 OSS
func UploadToQiniu(content string) (string, error) {
	accessKey := "BIWKkyBa1iJ-NEVPx0rSj4ol0fJBrF_F8rSPoiL8"
	secretKey := "OCUZEvmfyZ-_aThggjeAtCau_g2Mfzq8gPUnHYN7"
	bucket := "mybookdemo"
	domain := "http://t135by11b.hn-bkt.clouddn.com"
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	uploadToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan, // ✅ 华南 z2
		UseHTTPS:      false,               // 测试域名用 http
		UseCdnDomains: false,
	}
	formUploader := storage.NewFormUploader(&cfg)
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
	// 返回可访问 URL
	return fmt.Sprintf("%s/%s", domain, objectKey), nil
}

// 保存文章到 MongoDB
func SaveArticleToDB(article Article) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoURI := "mongodb://localhost:27017/book"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	collection := client.Database("book").Collection("articles")
	_, err = collection.InsertOne(ctx, article)
	return err
}
func main() {
	content := "这是文章内容。"
	title := "测试文章"
	authorId := "user123"
	url, err := UploadToQiniu(content)
	if err != nil {
		log.Fatal("上传 OSS 失败:", err)
	}
	fmt.Println("文章已上传，URL:", url)
	article := Article{
		Title:      title,
		AuthorId:   authorId,
		ContentURL: url,
		CreatedAt:  time.Now().UnixMilli(),
	}
	err = SaveArticleToDB(article)
	if err != nil {
		log.Fatal("保存文章到数据库失败:", err)
	}
	fmt.Println("文章信息已保存到数据库")
}

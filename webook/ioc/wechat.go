package ioc

import "github.com/hong-l1/project/webook/internal/service/oauth2/wechat"

func InitOauth2WechatService() wechat.WechatService {
	appid := "12341"
	appsecret := "12342"
	return wechat.NewWechatService(appid, appsecret)
}

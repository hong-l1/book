package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hong-l1/project/webook/internal/domain"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type WechatService interface {
	AuthURL(ctx context.Context, states string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
const targetPattern = `https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type WechatServiceImpl struct {
	appid     string
	appSecret string
	client    *http.Client
}

func NewWechatService(appid string, appSecret string) WechatService {
	return &WechatServiceImpl{
		appid:     appid,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}
func (w *WechatServiceImpl) AuthURL(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf(authURLPattern, w.appid, redirectURL, state), nil
}
func (w *WechatServiceImpl) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	targetURL := fmt.Sprintf(targetPattern, w.appid, redirectURL, code)
	rep, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := w.client.Do(rep)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&resp)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("调用微信接口失败 errcode %d, errmsg %s", res.ErrCode, res.ErrMsg)
	}
	zap.L().Info("调用微信拿到用户信息", zap.String("unionid", res.UnionID),
		zap.String("openid", res.opnenID))
	return domain.WechatInfo{
		OpenID:  res.opnenID,
		UnionId: res.UnionID,
	}, nil
}

type Result struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	UnionID      string `json:"unionid"`
	opnenID      string `json:"opnenid"`
}

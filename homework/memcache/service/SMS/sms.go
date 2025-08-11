package SMS

import (
	"encoding/json"
	"errors"
	sms "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type SMS struct {
	AccessKeyId     string //密钥id
	AccessKeySecret string //密钥
	SignName        string //签名
}

func NewSMS(accessKeyId string, accessKeySecret string, SignName string) *SMS {
	return &SMS{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		SignName:        SignName,
	}
}
func (s *SMS) SendSMS(template string, args []string, numbers string) error {
	client, err := sms.NewClientWithAccessKey("cn-hangzhou", s.AccessKeyId, s.AccessKeySecret)
	if err != nil {
		return err
	}
	res := sms.CreateSendSmsRequest()
	res.SignName = s.SignName
	res.PhoneNumbers = numbers
	param := map[string]string{
		"code": args[0],
	}
	bytes, err := json.Marshal(param)
	if err != nil {
		return err
	}
	res.TemplateParam = string(bytes)
	res.TemplateCode = template
	response, err := client.SendSms(res)
	if err != nil {
		return err
	}
	if response.Code == "OK" {
		return nil
	}
	return errors.New("发送失败")
}

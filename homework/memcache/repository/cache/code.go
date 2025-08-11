package cache

import (
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"strconv"
)

var (
	ErrCodeExhausted = errors.New("尝试次数已耗尽")
	ErrCodeSent      = errors.New("ErrCodeSent")
	ErrCodeInvalid   = errors.New("验证码错误")
)

type CodeCache struct {
	client *memcache.Client
}

func NewCodeCache(client *memcache.Client) *CodeCache {
	return &CodeCache{
		client: client,
	}
}

// key:业务+电话
func (c *CodeCache) SetCode(biz, phone, code string) error {
	key := biz + phone
	cntKey := key + ":cnt"
	err := c.client.Add(&memcache.Item{
		Key:        key,
		Value:      []byte(code),
		Expiration: 60, //60秒刷一次
	}) //key已经存在，说明已经发送过了，并且这个key没有过期
	if errors.Is(err, memcache.ErrNotStored) {
		return ErrCodeSent
	}
	if err != nil {
		return err
	}
	return c.client.Add(&memcache.Item{
		Key:        cntKey,
		Value:      []byte("3"),
		Expiration: 60,
	})
}

func (c *CodeCache) VerifyCode(biz, phone, inputCode string) error {
	//get不到过期时间
	key := biz + phone
	cntKey := key + ":cnt"
	res1, err := c.client.Get(key) //没有这个码
	if errors.Is(err, memcache.ErrCacheMiss) {
		return errors.New("ErrCodeNotSent")
	} else if err != nil {
		return err
	}
	res2, _ := c.client.Get(cntKey)
	count, err := strconv.Atoi(string(res2.Value))
	if err != nil {
		return errors.New("格式错误")
	}
	if count <= 0 {
		return ErrCodeExhausted
	}
	if string(res1.Value) != inputCode {
		count--
		res2.Value = []byte(strconv.Itoa(count))
		err = c.client.Set(res2) // 写回更新后的尝试次数
		if err != nil {
			return err
		}
		return ErrCodeInvalid
	} else {
		count = -1
		res2.Value = []byte(strconv.Itoa(count))
	}
	return nil
}

//service 负责生成验证码，调用repository,cache完成验证码的保存，然后调用短信业务发送验证码
//repository
//cache 完成对验证码的保存与验证

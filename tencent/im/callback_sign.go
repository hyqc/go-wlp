package tencent_im

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type ISign interface {
	Make(requestTime int64) string
	Check(sign string, requestTime int64) error
}

type Sign struct {
	options *SignOptions
}

type SignOptions struct {
	Token        string `json:"token"`         // 鉴权的token
	Expire       int64  `json:"expire"`        // 签名过期时间秒
	IgnoreExpire bool   `json:"ignore_expire"` // 是否忽略时间验证
}

func NewSign(options *SignOptions) *Sign {
	return &Sign{
		options: options,
	}
}

func (s *Sign) Make(requestTime int64) string {
	str := fmt.Sprintf("%s%d", s.options.Token, requestTime)
	h := sha256.New()
	h.Write([]byte(str))
	hb := h.Sum(nil)
	return hex.EncodeToString(hb)
}

func (s *Sign) Check(sign string, requestTime int64) error {
	_, err := s.CheckResult(sign, requestTime)
	return err
}

func (s *Sign) expire(requestTime int64) error {
	if s.options.IgnoreExpire {
		return nil
	}
	now := time.Now()
	after := CompareTime(requestTime, now.Unix(), time.Duration(s.options.Expire)*time.Second)
	if after {
		return nil
	}
	return fmt.Errorf("sign expired, request time: %d, now: %d", requestTime, now.Unix())
}

func (s *Sign) CheckResult(sign string, requestTime int64) (string, error) {
	if err := s.expire(requestTime); err != nil {
		return "", err
	}
	if sign == "" {
		return "", errors.New("target sign is empty")
	}
	ns := s.Make(requestTime)
	if sign == ns {
		return ns, nil
	}
	return ns, errors.New(fmt.Sprintf("check sign failed, remote sign: %s, new sign: %s", sign, ns))
}

// CompareTime 比较a+s秒后>b
func CompareTime(a, b int64, s time.Duration) bool {
	return time.Unix(a, 0).Add(s).After(time.Unix(b, 0))
}

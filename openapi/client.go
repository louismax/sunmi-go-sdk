package openapi

import (
	"errors"
)

type Client interface {
	Request(url string, params interface{}, headers map[string]string) ([]byte, error)
}

var (
	defaultParams = struct{}{}
	PublicKeyErr  = errors.New("公钥无效！")
	PrivateKeyErr = errors.New("私钥无效！")
	VerifySignErr = errors.New("签名校验失败！")
)

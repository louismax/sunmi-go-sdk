package openapi

import (
	"bytes"
	"crypto"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type RsaClient struct {
	Client
	AppId      string
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// Request 发送请求
func (c *RsaClient) Request(url string, params interface{}, headers map[string]string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			// Disable HTTP/2
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		},
	}
	if params == nil {
		params = defaultParams
	}
	bodyByte, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyByte))
	if err != nil {
		return nil, err
	}
	signHeader, err := c.SignRsa(string(bodyByte))
	if err != nil {
		return nil, err
	}
	for key, val := range signHeader {
		req.Header.Add(key, val)
	}
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	resp, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SignRsa rsa签名
func (c *RsaClient) SignRsa(data string) (map[string]string, error) {
	timestamp := time.Now().Unix()
	nonce := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	timestampStr := strconv.FormatInt(timestamp, 10)
	hash := sha256.New()
	hash.Write([]byte(data + c.AppId + timestampStr + nonce))
	signature, err := rsa.SignPKCS1v15(cryptoRand.Reader, c.PrivateKey, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"Sunmi-Timestamp": timestampStr,
		"Sunmi-Nonce":     nonce,
		"Sunmi-Appid":     c.AppId,
		"Sunmi-Sign":      base64.StdEncoding.EncodeToString(signature),
	}, nil
}

// VerifySignRsa rsa签名校验
func (c *RsaClient) VerifySignRsa(data, reqSign string) error {
	signByte, err := base64.StdEncoding.DecodeString(reqSign)
	shaNew := sha256.New()
	shaNew.Write([]byte(data))
	hashByte := shaNew.Sum(nil)
	err = rsa.VerifyPKCS1v15(c.PublicKey, crypto.SHA256, hashByte, signByte)
	if err != nil {
		return VerifySignErr
	}
	return nil
}

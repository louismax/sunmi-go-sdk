package openapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type HmacClient struct {
	Client
	AppId  string
	AppKey string
}

// Request 发送请求
func (c *HmacClient) Request(url string, params interface{}, headers map[string]string) ([]byte, error) {
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

	signHeader, err := c.SignHmac(c.AppId, string(bodyByte))
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

// SignHmac hmac签名
func (c *HmacClient) SignHmac(appId, data string) (map[string]string, error) {
	timestamp := time.Now().Unix()
	nonce := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	timestampStr := strconv.FormatInt(timestamp, 10)
	hash := hmac.New(sha256.New, []byte(c.AppKey)[:])
	hash.Write([]byte(data + appId + timestampStr + nonce))
	return map[string]string{
		"Sunmi-Timestamp": timestampStr,
		"Sunmi-Nonce":     nonce,
		"Sunmi-Appid":     appId,
		"Sunmi-Sign":      hex.EncodeToString(hash.Sum(nil)),
	}, nil
}

// VerifyHmac hmac签名校验
func (c *HmacClient) VerifyHmac(data, reqSign string) error {
	hashObj := hmac.New(sha256.New, []byte(c.AppKey)[:])
	hashObj.Write([]byte(data))
	sign := hex.EncodeToString(hashObj.Sum(nil))
	if reqSign != sign {
		return VerifySignErr
	}
	return nil
}

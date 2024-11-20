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

//const maxShowBodySize = 1024 * 100
//
//type HttpClient struct {
//	Client  *resty.Client
//	Request *resty.Request
//
//	disableLog               bool            // default: false 默认打印日志(配置SetLog后)
//	disableMetrics           bool            // default: false 默认开启统计
//	disableBreaker           bool            // default: true 默认关闭熔断
//	slowThresholdMs          int64           // default: 0 默认关闭慢请求打印
//	hideRespBodyLogsWithPath map[string]bool // 不打印path在map里的返回体
//	hideReqBodyLogsWithPath  map[string]bool // 不打印path在map里的请求体
//	maxShowBodySize          int64
//}
//
//func NewHttpClient() *HttpClient {
//	// Create a Resty Client
//	client := resty.New()
//
//	// Retries are configured per client
//	client.
//		// Set retry count to non zero to enable retries
//		SetRetryCount(3).
//		// TimeOut
//		SetTimeout(5 * time.Second).
//		// You can override initial retry wait time.
//		// Default is 100 milliseconds.
//		SetRetryWaitTime(2 * time.Second).
//		// MaxWaitTime can be overridden as well.
//		// Default is 2 seconds.
//		SetRetryMaxWaitTime(5 * time.Second).
//		// SetRetryAfter sets callback to calculate wait time between retries.
//		// Default (nil) implies exponential backoff with jitter
//		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
//			return 0, errors.New("quota exceeded")
//		})
//
//	return &HttpClient{
//		Client:          client,
//		Request:         client.R(),
//		disableMetrics:  false,
//		disableLog:      false,
//		disableBreaker:  true, // default disable, will open soon
//		maxShowBodySize: maxShowBodySize,
//	}
//}

# sunmi-go-sdk
[简体中文](README.md)
## SunMI商米开放对接API，golang-SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/louismax/sunmi-go-sdk)](https://goreportcard.com/report/github.com/louismax/sunmi-go-sdk)
[![GoDoc](https://godoc.org/github.com/louismax/sunmi-go-sdk?status.svg)](https://godoc.org/github.com/louismax/sunmi-go-sdk)
[![GitHub release](https://img.shields.io/github/tag/louismax/sunmi-go-sdk.svg)](https://github.com/louismax/sunmi-go-sdk/releases)
[![GitHub license](https://img.shields.io/github/license/louismax/sunmi-go-sdk.svg)](https://github.com/louismax/sunmi-go-sdk/blob/master/LICENSE)
[![GitHub Repo Size](https://img.shields.io/github/repo-size/louismax/sunmi-go-sdk.svg)](https://img.shields.io/github/repo-size/louismax/sunmi-go-sdk.svg)
[![GitHub Last Commit](https://img.shields.io/github/last-commit/louismax/sunmi-go-sdk.svg)](https://img.shields.io/github/last-commit/louismax/sunmi-go-sdk.svg)

## 安装
`go get -v github.com/louismax/sunmi-go-sdk`

## 使用
### 获取OpenApi_client(hmac)
```
client := sunmi_go_sdk.NewHmacClient("appid", "appKey")
```
### OpenApiClient实现的方法

#### client.DeviceBindShop  打印机绑定店铺
#### client.DeviceUnbindShop  打印机解绑店铺
#### client.DeviceOnlineStatus  查询渠道下所有设备或者单台设备在线状态
#### client.DeviceClearPrintJob  清除云端缓存的打印队列
#### client.DevicePushVoice  给打印机推送播报语音
#### client.TicketPrintStatus  查询订单打印状态
#### client.DevicePushContent  给打印机推送订单详情(直推模式)

### OpenApi使用示例
```
//查询打印机在线状态
resp, err := client.DeviceOnlineStatus("sn", 1, 10)
if err != nil {
	t.Log(err)
}
t.Log(string(resp))
```



## 参考资料
* [Sunmi开发者文档中心](https://developer.sunmi.com/docs/zh-CN/index)
* [云打印机开发对接API](https://developer.sunmi.com/docs/zh-CN/xeghjk491/fmqeghjk513)
* [云打印机指令集](https://developer.sunmi.com/docs/zh-CN/xeghjk491/fzqeghjk513)

## 协议
MIT 许可证（MIT）。有关更多信息，请参见[协议文件](LICENSE)。
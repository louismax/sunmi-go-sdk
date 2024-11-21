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
### 🚀 获取OpenApi_client(hmac)
```go
client := sunmi_go_sdk.NewHmacClient("appid", "appKey")
```

### 🚀 OpenApiClient实现的方法

##### client.DeviceBindShop  *打印机绑定店铺*
##### client.DeviceUnbindShop  *打印机解绑店铺*
##### client.DeviceOnlineStatus  *查询渠道下所有设备或者单台设备在线状态*
##### client.DeviceClearPrintJob  *清除云端缓存的打印队列*
##### client.DevicePushVoice  *给打印机推送播报语音*
##### client.TicketPrintStatus  *查询订单打印状态*
##### client.DevicePushContent  *给打印机推送订单详情(直推模式)*

### 🚀 OpenApi使用示例

```go
//查询打印机在线状态
resp, err := client.DeviceOnlineStatus("sn", 1, 10)
if err != nil {
	t.Log(err)
}
t.Log(string(resp))
```
### 🚀 打印指令集


##### NewPrint *打印对象构造函数*
```go
//打印对象构造函数 NewPrint
p := printer.NewPrint()//NewPrint支持提供可变参数,不传默认为384
//或
p := printer.NewPrint(384) //NewPrint支持提供可变参数,不传默认为384
```
✨常见的打印机分辨率有203dpi，300dpi，600dpi等。市场上的打印机以203dpi为主。dpi是每英寸的点数，203dpi 也就是203点/inch，1inch=25.4mm，也就是说 203点/25.4mm， 换算下来就是 8 点/mm。
✨热敏打印机都是以点为单位来进行排版计算打印机的有效打印宽度和纸张宽度是不一样的。常见的小票机规格书上描述的宽度如 58mm、80mm等。

✨纸张宽度=打印宽度+左右边距
| 纸张宽度 | 有效打印宽度 |  页边距 |
|--|--|--|
| 58mm | 48mm | 5mm |
| 80mm | 72mm | 4mm |
✨如果机器是203dpi，那么对应的有效打印点数就是
| 有效打印宽度 | 有效打印点数 |  
|--|--|
| 48mm | 384点（48*8） |
| 72mm | 576点（72*8） |

##### *常用打印对象实现方法*
```go
p.AppendText("string") //添加文字
p.AppendDivider() //分割线
p.AppendDivider("分割线") //分割线，AppendDivider(...string)提供可变参数指定分割线居中文字
p.LineFeed() //[LF]打印缓冲区和进纸行中的数据
p.LineFeed(2) //[LF]打印缓冲区和进纸行中的数据,LineFeed(...int)提供可变参数指定执行次数
p.RestoreDefaultSettings() //[ESC @] 恢复默认设置
p.RestoreDefaultLineSpacing() //[ESC 2] 恢复默认行距
p.SetLineSpacing(3) //[ESC 3] 设置行距
p.SetPrintModes(true,false,true) //[ESC !] 设置打印模式,参数指定加粗、倍高、倍宽
p.SetCharacterSize(1,1) //[GS !] 设置字符大小，参数指定高、宽
p.HorizontalTab() //[HT] 插入水平制表符，HorizontalTab(...int)提供可变参数指定执行次数
p.SetAbsolutePrintPosition(1) //[ESC $] 设置绝对打印位置.
p.SetRelativePrintPosition(1) //[ESC \] 设置相对打印位置.
p.SetAlignment(printer.AlignCenter) //[ESC a] 设置对齐方式.
p.SetUnderlineMode(1) //[ESC -] 设置下划线模式.
p.SetBlackWhiteReverseMode(true) //[GS B] 设置/取消黑白倒转模式.
p.SetUpsideDownMode(true) //[ESC {] 设置/取消倒立模式.
p.SetBold(true) // [ESC E] 设置/取消倒立模式.
p.CutPaper(true) //[GS V m] 切纸.
p.PostponedCutPaper(true,1) // [GS V m n] 延期裁纸. 在收到此命令后，打印机将不执行切割，直到(d + n)点线馈送，其中d是打印位置和切割位置之间的距离
p.SetupColumns([][]int{{210, printer.AlignLeft, printer.ColumnFlagBold}, {0, printer.AlignRight, printer.ColumnFlagBold}}) // 设置列
p.PrintInColumns([]string{"支付时间:", "2024-09-08 00:00:35"}) //按列打印
p.AppendBarcode(printer.HRI_POS_ABOVE, 68, 2, 73, "243110163822867492208888") // 添加条形码
p.AppendQRCode(8, 2, "https://github.com/louismax/sunmi-go-sdk") //添加二维码
......
```

##### 使用示例

```go
func TestDoPrint(t *testing.T) {
    //获取HmacClient
    client := NewHmacClient("appid", "appKey")  
    //查询云打印机设备在线情况
    resp, err := client.DeviceOnlineStatus("sn", 1, 10)  
    if err != nil {  
       t.Log(err)  
    }  
    t.Log(string(resp))  
    
    //打印58mm小票示例
    p := printer.NewPrint()  
    p.SetAlignment(printer.AlignCenter)  
    p.SetCharacterSize(2, 2)  
    p.SetBold(true)  
    p.AppendText("XX水果店")  
    p.SetBold(false)  
    p.LineFeed(2)  
    p.SetCharacterSize(2, 1)  
    p.SetAlignment(printer.AlignLeft)  
    p.AppendText("收货地址:XX小区-X栋X单元1001")  
    p.SetBold(false)  
    p.LineFeed()  
    p.AppendText("收货人:张(先生)")  
    p.LineFeed()  
    p.AppendText("电话:131****0001")  
    p.LineFeed()  
    p.AppendDivider("明细")  
    p.SetAlignment(printer.AlignLeft)  
    p.SetupColumns([][]int{{210, printer.AlignLeft, printer.ColumnFlagBold}, {70, printer.AlignCenter, printer.ColumnFlagBold}, {0, printer.AlignRight, printer.ColumnFlagBold}})  
    p.PrintInColumns([]string{"商品名称", "数量", "小计(元)"})  
    p.SetupColumns([][]int{{210, printer.AlignLeft, 0}, {70, printer.AlignCenter, 0}, {0, printer.AlignRight, printer.ColumnFlagBold}})  
    p.PrintInColumns([]string{"西瓜", "x1", "￥20.99"})  
    p.PrintInColumns([]string{"苹果", "x1000", "￥20.99"})   
    p.AppendDivider()  
    p.SetupColumns([][]int{{110, printer.AlignLeft, 0}, {0, printer.AlignRight, printer.ColumnFlagBold}})  
    p.PrintInColumns([]string{"商品总数:", "X2"})  
    p.PrintInColumns([]string{"实付金额:", "￥198.53"})  
    p.AppendDivider()  
    p.SetupColumns([][]int{{110, printer.AlignLeft, 0}, {0, printer.AlignRight, 0}})  
    p.PrintInColumns([]string{"订单编号:", "242520278490929129688888"})  
    p.PrintInColumns([]string{"下单时间:", "2024-09-08 00:00:35"})  
    p.PrintInColumns([]string{"支付方式:", "微信在线支付"})  
    p.PrintInColumns([]string{"支付时间:", "2024-09-08 00:00:35"})  
    p.AppendDivider()  
    p.SetAlignment(printer.AlignCenter)  
    p.AppendText("微信扫描下方二维码，可在小程序内查看订单详情!")  
    p.LineFeed()  
    p.AppendQRCode(8, 2, "http://natfrp.louiss.net")  
    p.LineFeed(2)  
    p.AppendDivider("裁--纸--线")  
    p.LineFeed(2)

    //调用直推打印API
    resp2, err := client.DevicePushContent("sn", kTool.MakeYearDaysRand(), p.Content, "新订单语音播报内容", "", 1, 1, 1)  
    if err != nil {  
       t.Log(err)  
    }  
    t.Log(string(resp2))  
}
```

## 参考资料
* [Sunmi开发者文档中心](https://developer.sunmi.com/docs/zh-CN/index)
* [云打印机开发对接API](https://developer.sunmi.com/docs/zh-CN/xeghjk491/fmqeghjk513)
* [云打印机指令集](https://developer.sunmi.com/docs/zh-CN/xeghjk491/fzqeghjk513)

## 协议
Apache-License2.0。有关更多信息，请参见[协议文件](LICENSE)。


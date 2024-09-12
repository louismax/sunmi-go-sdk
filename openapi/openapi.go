package openapi

const OPENAPIURL string = "https://openapi.sunmi.com"

const (
	DeviceBindShopUrl      string = "/v2/printer/open/open/device/bindShop"
	DeviceUnbindShopUrl    string = "/v2/printer/open/open/device/unbindShop"
	DeviceOnlineStatusUrl  string = "/v2/printer/open/open/device/onlineStatus"
	DeviceClearPrintJobUrl string = "/v2/printer/open/open/device/clearPrintJob"
	DevicePushVoiceUrl     string = "/v2/printer/open/open/device/pushVoice"
	TicketPrintStatusUrl   string = "/v2/printer/open/open/ticket/printStatus"
	DevicePushContentUrl   string = "//v2/printer/open/open/device/pushContent"
)

// DeviceBindShop 打印机绑定店铺
func (c *HmacClient) DeviceBindShop(sn string, shopId int) ([]byte, error) {
	return c.Request(OPENAPIURL+DeviceBindShopUrl, map[string]interface{}{
		"sn":      sn,
		"shop_id": shopId,
	}, map[string]string{
		"Source":       "openapi",
		"Content-Type": "application/json",
	})
}

// DeviceUnbindShop 打印机解绑店铺
func (c *HmacClient) DeviceUnbindShop(sn string, shopId int) ([]byte, error) {
	return c.Request(OPENAPIURL+DeviceUnbindShopUrl, map[string]interface{}{
		"sn":      sn,
		"shop_id": shopId,
	}, map[string]string{
		"Source":       "openapi",
		"Content-Type": "application/json",
	})
}

// DeviceOnlineStatus 查询渠道下所有设备或者单台设备在线状态
func (c *HmacClient) DeviceOnlineStatus(sn string, pageNo, pageSize int) ([]byte, error) {
	req := map[string]interface{}{}
	if sn != "" {
		req["sn"] = sn
	}
	if pageNo > 0 && pageSize > 0 {
		req["page_no"] = pageNo
		req["page_size"] = pageSize
	}
	return c.Request(OPENAPIURL+DeviceOnlineStatusUrl, req, map[string]string{
		//"Connection":   "no-cache",
		"Source":       "openapi",
		"Content-Type": "application/json",
	})

}

// DeviceClearPrintJob 清除云端缓存的打印队列
func (c *HmacClient) DeviceClearPrintJob(sn string) ([]byte, error) {
	return c.Request(OPENAPIURL+DeviceClearPrintJobUrl, map[string]interface{}{
		"sn": sn,
	}, map[string]string{
		"Source":       "openapi",
		"Content-Type": "application/json",
	})
}

// DevicePushVoice 给打印机推送播报语音
func (c *HmacClient) DevicePushVoice(sn, content, mediaUrl string, expireIn, cycle, interval int) ([]byte, error) {
	req := map[string]interface{}{
		"sn": sn,
	}

	if content != "" {
		req["content"] = content
	} else if mediaUrl != "" {
		req["media_url"] = mediaUrl
	}

	if expireIn > 0 {
		req["expire_in"] = expireIn
	}
	if cycle > 0 {
		req["cycle"] = cycle
	}
	if interval > 0 {
		req["interval"] = interval
	}

	return c.Request(OPENAPIURL+DevicePushVoiceUrl, req, map[string]string{
		"Source":       "openapi",
		"Content-Type": "application/json",
	})
}

// TicketPrintStatus 查询订单打印状态
func (c *HmacClient) TicketPrintStatus(tradeNo string) ([]byte, error) {
	return c.Request(OPENAPIURL+TicketPrintStatusUrl, map[string]interface{}{
		"trade_no": tradeNo,
	}, map[string]string{
		"Source":       "openapi",
		"Content-Type": "application/json",
	})
}

// DevicePushContent 给打印机推送订单详情(直推模式)
func (c *HmacClient) DevicePushContent(sn, tradeNo, content, mediaText, mediaUrl string, count, orderType, cycle int) ([]byte, error) {
	req := map[string]interface{}{
		"trade_no": tradeNo,
		"sn":       sn,
		"content":  content,
		"count":    count,
	}

	if mediaText != "" {
		req["media_text"] = mediaText
	} else if mediaUrl != "" {
		req["media_url"] = mediaUrl
	}

	if orderType > 0 {
		req["order_type"] = orderType
	}

	if cycle > 0 {
		req["cycle"] = cycle
	}

	return c.Request(OPENAPIURL+DevicePushContentUrl, req, map[string]string{
		"Source":       "openapi",
		"Content-Type": "application/json",
	})
}

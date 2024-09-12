package sunmi_go_sdk

import (
	"testing"
)

func TestDeviceOnlineStatus(t *testing.T) {
	client := NewHmacClient("appid", "appKey")
	resp, err := client.DeviceOnlineStatus("sn", 1, 10)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(resp))
}

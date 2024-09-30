package sunmi_go_sdk

import "testing"

func TestDOSOnline(t *testing.T) {
	client := NewHmacClient("cf2d989f2b4a40d7b5ca2e23c3d3fd8c", "f9ab393ef7bd444b8ba2cc3e9351cd43")
	resp, err := client.DeviceOnlineStatus("N302236R40613", 1, 10)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(resp))

}

func TestDOSbb(t *testing.T) {
	client := NewHmacClient("cf2d989f2b4a40d7b5ca2e23c3d3fd8c", "f9ab393ef7bd444b8ba2cc3e9351cd43")
	resp, err := client.DeviceOnlineStatus("N302236R40613", 1, 10)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(resp))

	_, _ = client.DevicePushVoice("N302236R40613", "技术员测试是否在线，手动推送", "", 0, 0, 0)

}

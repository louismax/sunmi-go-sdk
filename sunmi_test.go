package sunmi_go_sdk

import (
	"github.com/louismax/sunmi-go-sdk/printer"
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

func TestPrint(t *testing.T) {
	p := printer.PrintObject{}

	p.SetCharacterSize(8, 8)

	t.Log(p.Content)
}

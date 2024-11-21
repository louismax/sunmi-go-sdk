package sunmi_go_sdk

import (
	"github.com/louismax/sunmi-go-sdk/printer"
	"github.com/louismax/sunmi-go-sdk/tools"
	"testing"
)

func TestDeviceOnlineStatus(t *testing.T) {
	client := NewHmacClient("appid", "appKey")
	resp, err := client.DeviceOnlineStatus("sn", 1, 10)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(resp))

	client2, err := NewRsaClient("appid", "privateKey", "publicKey")
	if err != nil {
		t.Log(err)
	}
	_, _ = client2.SignRsa("data")
}

func TestPrint(t *testing.T) {
	p := printer.NewPrint()

	p.SetCharacterSize(8, 8)

	t.Log(p.Content)
}

func TestGetRandomString(t *testing.T) {
	t.Log(tools.GetRandomString(10))
}

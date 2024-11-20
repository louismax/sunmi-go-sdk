package sunmi_go_sdk

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/louismax/sunmi-go-sdk/openapi"
)

// NewHmacClient return a hmac client
func NewHmacClient(appId, appKey string) *openapi.HmacClient {
	return &openapi.HmacClient{
		AppId:  appId,
		AppKey: appKey,
	}
}

// NewRsaClient return a rsa client
func NewRsaClient(appId, privateKey, publicKey string) (*openapi.RsaClient, error) {
	prvBlock, _ := pem.Decode([]byte(privateKey))
	if prvBlock == nil {
		return nil, openapi.PrivateKeyErr
	}

	prv, err := x509.ParsePKCS8PrivateKey(prvBlock.Bytes)
	if err != nil {
		return nil, openapi.PrivateKeyErr
	}

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, openapi.PublicKeyErr
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, openapi.PublicKeyErr
	}
	pub := pubInterface.(*rsa.PublicKey)
	return &openapi.RsaClient{
		AppId:      appId,
		PublicKey:  pub,
		PrivateKey: prv.(*rsa.PrivateKey),
	}, nil
}

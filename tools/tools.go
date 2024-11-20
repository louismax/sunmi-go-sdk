package tools

import (
	cryptorand "crypto/rand"
	"fmt"
)

// GetRandomString 获取指定长度的随机字符串
func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	_, _ = cryptorand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

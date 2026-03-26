package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

// Base64Encode 对数据进行 Base64 编码
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64EncodeByte 对数据进行 Base64 编码
func Base64EncodeByte(data []byte) []byte {
	return []byte(Base64Encode(data))
}

// Base64Decode 对 Base64 编码的数据进行解码
func Base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

package sms

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// PKCS5Padding 对明文进行填充以符合块大小
func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// AESEncryptECB 实现 AES/ECB/PKCS5Padding 加密
func AESEncryptECB(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(src)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("plaintext is not a multiple of the block size")
	}

	encrypted := make([]byte, len(src))
	for bs, be := 0, block.BlockSize(); bs < len(src); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(encrypted[bs:be], src[bs:be])
	}

	return encrypted, nil
}

// GenerateSignature 根据 appid、apiname、timestamp 和 appkey 生成签名
func GenerateSignature(appid, apiname, appkey string) (string, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix()) // 10位时间戳
	stringToSign := appid + apiname + timestamp
	key := strings.ReplaceAll(appkey, "-", "")
	keyBytes := []byte(key)

	// key 长度必须是 16, 24 或 32 字节（AES 的要求）
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", fmt.Errorf("invalid key length: %d (must be 16, 24, or 32 bytes)", len(keyBytes))
	}

	paddedText := PKCS5Padding([]byte(stringToSign), aes.BlockSize)
	encrypted, err := AESEncryptECB(paddedText, keyBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}


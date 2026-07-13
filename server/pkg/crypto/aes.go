package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// 默认密钥（32 字节用于 AES-256），生产环境应通过配置文件覆盖
var defaultKey = []byte("monitor-web-change-default-key32")

var encKey []byte

func SetKey(key string) {
	if len(key) >= 32 {
		encKey = []byte(key[:32])
	} else if len(key) > 0 {
		padded := make([]byte, 32)
		copy(padded, key)
		encKey = padded
	} else {
		encKey = defaultKey
	}
}

func getKey() []byte {
	if len(encKey) == 0 {
		return defaultKey
	}
	return encKey
}

// Encrypt AES-256-GCM 加密，返回 Base64 编码的密文
func Encrypt(plaintext string) (string, error) {
	key := getKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt AES-256-GCM 解密 Base64 编码的密文
func Decrypt(encoded string) (string, error) {
	key := getKey()
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

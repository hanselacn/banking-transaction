package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func EncryptAES(dataStr string, keyStr string) (string, error) {
	key := []byte(keyStr)
	data := []byte(dataStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], data)

	base64ciphertext := base64.StdEncoding.EncodeToString(ciphertext)

	return base64ciphertext, nil
}

func DecryptAES(base64ciphertext string, keyStr string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(base64ciphertext)
	if err != nil {
		return "", err
	}

	key := []byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

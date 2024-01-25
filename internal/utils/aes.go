package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func EncryptAES(key []byte, plaintext string) string {
	c, _ := aes.NewCipher(key)
	plaintextBytes := []byte(plaintext)
	paddedPlaintext := PKCS7Padding(plaintextBytes, aes.BlockSize)
	out := make([]byte, aes.BlockSize+len(paddedPlaintext))
	iv := out[:aes.BlockSize]
	c.Encrypt(iv, iv)
	mode := cipher.NewCBCEncrypter(c, iv)
	mode.CryptBlocks(out[aes.BlockSize:], paddedPlaintext)
	return hex.EncodeToString(out)
}

func DecryptAES(key []byte, ct string) string {
	ciphertext, _ := hex.DecodeString(ct)
	c, _ := aes.NewCipher(key)
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(c, iv)
	pt := make([]byte, len(ciphertext))
	mode.CryptBlocks(pt, ciphertext)
	unpaddedPlaintext := PKCS7Unpadding(pt)
	return string(unpaddedPlaintext)
}

func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func PKCS7Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

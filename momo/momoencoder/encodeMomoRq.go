package momoencoder

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func EncodeMomoRq(plaintext string, key string) (string, error) {
	blockSize := 16
	iv := ""
	for i := 0; i < blockSize; i++ {
		iv += string(rune(0))
	}
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := _PKCS5Padding([]byte(plaintext), blockSize, len(plaintext))
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func _PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

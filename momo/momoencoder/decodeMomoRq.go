package momoencoder

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func DecodeMomoRq(_cipherText string, encKey string) (string, error) {
	blockSize := 16
	iv := ""
	for i := 0; i < blockSize; i++ {
		iv += string(rune(0))
	}

	ciphertext, err := base64.StdEncoding.DecodeString(_cipherText)
	if err != nil {
		return "", err
	}

	bKey := []byte(encKey)
	bIV := []byte(iv)

	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks(ciphertext, ciphertext)
	return string(PKCS5Trimming(ciphertext)), nil
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

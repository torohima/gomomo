package momoencoder

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func DecodeRSA(_ciphertext string, _privKey string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(_ciphertext)
	if err != nil {
		return "", err
	}
	privKey := []byte(_privKey)
	block, _ := pem.Decode(privKey)
	if block == nil {
		return "", errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	dec, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return "", err
	}
	return string(dec), nil
}

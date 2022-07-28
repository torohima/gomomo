package momoencoder

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func EncodeRSA(_ciphertext string, _pubKey string) (string, error) {
	ciphertext := []byte(_ciphertext)
	pubKey := []byte(_pubKey)
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return "", errors.New("public key error!")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	enc, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), ciphertext)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enc), nil
}

func EncodeRSAMomoPubKey(_ciphertext string, _pubKey string) (string, error) {
	ciphertext := []byte(_ciphertext)
	pubKey := []byte(_pubKey)
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return "", errors.New("public key error!")
	}
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	enc, err := rsa.EncryptPKCS1v15(rand.Reader, pub, ciphertext)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(enc), nil
}

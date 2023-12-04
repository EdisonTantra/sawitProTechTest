package locker

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Locker struct {
	secretKey string
}

func New(secretKey string) *Locker {
	return &Locker{
		secretKey: secretKey,
	}
}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func (l *Locker) Decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (l *Locker) Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func (l *Locker) Encrypt(s string) (string, error) {
	block, err := aes.NewCipher([]byte(l.secretKey))
	if err != nil {
		return "", err
	}

	plainText := []byte(s)
	cipherText := make([]byte, len(plainText))
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cfb.XORKeyStream(cipherText, plainText)
	return l.Encode(cipherText), nil
}

func (l *Locker) Decrypt(base64Str string) (string, error) {
	block, err := aes.NewCipher([]byte(l.secretKey))
	if err != nil {
		return "", err
	}

	cipherText, err := l.Decode(base64Str)
	if err != nil {
		return "", err
	}

	plainText := make([]byte, len(cipherText))
	cfb := cipher.NewCFBDecrypter(block, bytes)
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

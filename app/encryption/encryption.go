package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/typomedia/patchouli/app"
	"io"
)

func EncryptString(stringToEncrypt string) (string, error) {
	key, err := app.GetApp().Config.GetCipherKey()
	if err != nil {
		return "", err
	}
	aesKey, _ := hex.DecodeString(key)
	plaintext := []byte(stringToEncrypt)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func DecryptString(encString string) (string, error) {
	key, err := app.GetApp().Config.GetCipherKey()

	aesKey, _ := hex.DecodeString(key)
	enc, _ := hex.DecodeString(encString)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

/*
Copyright © 2024 Yannick Foeillet <bzhtux@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package endefi

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	"log"

	"github.com/bzhtux/endefi/config"
)

type secretService struct {
	secretRepo SecretRepository
}

func (ss *secretService) GetSecretKey(cfg *config.Config) (*Secret, error) {
	return ss.secretRepo.GetSecretKey(cfg)
}

func NewSecretService(secretRepo SecretRepository) SecretService {
	return &secretService{
		secretRepo,
	}
}

// CheckKeySize check the key size.
func checkKeySize(keyString []byte) error {
	_, err := aes.NewCipher(keyString)
	if err != nil {
		log.Default().Printf("Error: %s", err.Error())
		return err
	}
	return nil
}

// GenerateRandomKey generate a random key
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// ConvertByteToString convert an array of byte to a string
func ConvertByteToString(b []byte) string {
	keyStr := hex.EncodeToString(b)

	return keyStr
}

// ConvertStringToByte convert a string to an array of byte
func ConvertStringToByte(s string) ([]byte, error) {
	b, err := hex.DecodeString(s)

	if err != nil {
		return nil, err
	}

	return b, nil
}

// B64Encode encode a string into base64
func B64Encode(s string) string {
	b64 := base64.URLEncoding.EncodeToString([]byte(s))

	return b64
}

// B64Decode decode a base64 string into an array of byte
func B64Decode(s string) ([]byte, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// EncryptData encrypt data using aes GCM
func EncryptData(plaindata string, key []byte) ([]byte, error) {
	if err := checkKeySize(key); err != nil {
		return nil, err
	}

	d := []byte(plaindata)
	k := []byte(key)

	c, err := aes.NewCipher(k)

	if err != nil {
		log.Default().Printf("Error: %s", err.Error())
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		log.Default().Printf("Error: %s", err.Error())
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Default().Printf("Error: %s", err.Error())
		return nil, err
	}

	encrypted := gcm.Seal(nonce, nonce, d, nil)

	return encrypted, nil
}

// DecryptData decrypt data using aes GCM
func DecryptData(cipherdata []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := cipherdata[:nonceSize], cipherdata[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

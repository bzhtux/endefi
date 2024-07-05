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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"

	// "encoding/hex"
	"io"
	// "io/ioutil"
	"os"

	"golang.org/x/crypto/pbkdf2"

	"github.com/bzhtux/endefi/config"
)

type Secret struct {
	Key      string `env:"ENDEFI_SECRET_KEY"`
	Provider string `env:"ENDEFI_SECRET_PROVIDER"`
}

type SecretRepository interface {
	GetSecretKey(*config.Config) (*Secret, error)
}

func EncryptSecretFile(source string, password []byte) error {

	if _, err := os.Stat(source); os.IsNotExist(err) {
		return err
	}

	// Read the file
	plaintext, err := os.ReadFile(source)
	if err != nil {
		panic(err.Error())
	}

	key := password
	nonce := make([]byte, 12)

	// Randomizing the nonce
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	dk := pbkdf2.Key(key, nonce, 4096, 32, sha1.New)

	block, err := aes.NewCipher(dk)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// Append the nonce to the end of file
	ciphertext = append(ciphertext, nonce...)

	f, err := os.Create(source)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.Copy(f, bytes.NewReader(ciphertext))
	if err != nil {
		panic(err.Error())
	}

	return nil
}

// func DecryptSecretFile() {

// }

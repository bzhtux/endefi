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
package config

import (
	"errors"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
)

const (
	ENV_PREFIX = "ENDEFI"
)

type Config struct {
	App    *AppConfig
	Secret *SecretConfig
}

type AppConfig struct {
	Name     string `env:"ENDEFI_APP_NAME"`
	Provider string `env:"ENDEFI_SECRET_PROVIDER"`
}

type SecretConfig struct {
	Name         string `env:"ENDEFI_SECRET_NAME"`
	File         string `env:"ENDEFI_SECRET_FILE"`
	Key          string `env:"ENDEFI_SECRET_KEY"`
	URL          string `env:"ENDEFI_SECRET_URL"`
	Username     string `env:"ENDEFI_SECRET_USERNAME"`
	Password     string `env:"ENDEFI_SECRET_PASSWORD"`
	ClientID     string `env:"ENDEFI_SECRET_CLIENT_ID"`
	ClientSecret string `env:"ENDEFI_SECRET_CLIENT_SECRET"`
}

func setAppDefault() *AppConfig {
	// set default config
	return &AppConfig{
		Name: "EnDeFi",
	}
}

func setSecretDefault() *SecretConfig {
	// set default config
	homedir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return &SecretConfig{
		File: homedir + "/.endefi/secret.yaml",
	}
}

func NewAppConfig() (*AppConfig, error) {
	ac := setAppDefault()
	if eap := os.Getenv(ENV_PREFIX + "_SECRET_PROVIDER"); eap != "" {
		ac.Provider = eap
	} else {
		return nil, errors.New("no secret provider found via environment variable")
	}
	if ean := os.Getenv(ENV_PREFIX + "_APP_NAME"); ean != "" {
		ac.Name = ean
	}
	return ac, nil
}

func NewSecretConfig() (*SecretConfig, error) {
	sc := setSecretDefault()
	switch os.Getenv(ENV_PREFIX + "_SECRET_PROVIDER") {
	case "env":
		if esk := os.Getenv(ENV_PREFIX + "_SECRET_KEY"); esk != "" {
			sc.Key = esk
			sc.Name = "env"
			sc.File = ""
			return sc, nil
		} else {
			return nil, errors.New("no secret key found via environment variable")
		}
	case "local":
		if esf := os.Getenv(ENV_PREFIX + "_SECRET_FILE"); esf != "" {
			sc.File = esf
			sc.Name = "local"
		}
		return sc, nil
	default:
		return nil, errors.New("no secret provider found via environment variable")
	}
}

func fileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

func NewConfig() (*Config, error) {
	ac, err := NewAppConfig()
	if err != nil {
		return nil, err
	}
	sc, err := NewSecretConfig()
	if err != nil {
		return nil, err
	}
	if !fileExists(sc.File) {
		return nil, errors.New("secret file not found")
	}
	return &Config{
		App:    ac,
		Secret: sc,
	}, nil
}

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
		Name:     "EnDeFi",
		Provider: "env",
	}
}

func setSecretDefault() *SecretConfig {
	// set default config
	homedir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return &SecretConfig{
		Name: "env",
		File: homedir + "/.endefi/secret.yaml",
	}
}

func NewAppConfig() *AppConfig {
	ac := setAppDefault()
	if ean := os.Getenv(ENV_PREFIX + "_APP_NAME"); ean != "" {
		ac.Name = ean
	}
	if eap := os.Getenv(ENV_PREFIX + "_SECRET_PROVIDER"); eap != "" {
		ac.Provider = eap
	}
	return ac
}

func NewSecretConfig() *SecretConfig {
	sc := setSecretDefault()
	switch os.Getenv(ENV_PREFIX + "_SECRET_PROVIDER") {
	case "env":
		if esk := os.Getenv(ENV_PREFIX + "_SECRET_KEY"); esk != "" {
			sc.Key = esk
		}
		return sc
	case "bitwarden":
		// Do something
		// if esn := os.Getenv(ENV_PREFIX + "_SECRET_NAME"); esn != "" {
		// 	sc.Name = esn
		// }

		// if esp := os.Getenv(ENV_PREFIX + "_SECRET_PROVIDER"); esp != "" {
		// 	sc.Provider = esp
		// }
		// if esci := os.Getenv(ENV_PREFIX + "_SECRET_CLIENT_ID"); esci != "" {
		// 	sc.ClientID = esci
		// }
		// if escs := os.Getenv(ENV_PREFIX + "_SECRET_CLIENT_SECRET"); escs != "" {
		// 	sc.ClientSecret = escs
		// }
		// if esurl := os.Getenv(ENV_PREFIX + "_SECRET_URL"); esurl != "" {
		// 	sc.URL = esurl
		// }
		// if esu := os.Getenv(ENV_PREFIX + "_SECRET_USERNAME"); esu != "" {
		// 	sc.Username = esu
		// }
		// if espwd := os.Getenv(ENV_PREFIX + "_SECRET_PASSWORD"); espwd != "" {
		// 	sc.Password = espwd
		// }
		return sc
	case "local":
		if esf := os.Getenv(ENV_PREFIX + "_SECRET_FILE"); esf != "" {
			sc.File = esf
		}
		return sc
	default:
		log.Fatal("No secret provider found, exiting ...")
		os.Exit(1)
	}

	return sc
}

func NewConfig() *Config {
	ac := NewAppConfig()
	sc := NewSecretConfig()
	return &Config{
		App:    ac,
		Secret: sc,
	}
}

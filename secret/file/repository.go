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
package file

import (
	"os"

	"github.com/bzhtux/endefi/config"
	"github.com/bzhtux/endefi/endefi"

	"gopkg.in/yaml.v3"
)

type FileRepository struct {
	cfg *config.SecretConfig
}

func NewFileRepository(cfg config.Config) endefi.SecretRepository {
	return &FileRepository{
		cfg: cfg.Secret,
	}
}

func (r *FileRepository) GetSecretKey(cfg *config.Config) (*endefi.Secret, error) {
	sk, err := config.NewSecretConfig()
	if err != nil {
		return nil, err
	}
	type cfgYML struct {
		Key      string `yaml:"key"`
		Provider string `yaml:"provider"`
	}
	cy := cfgYML{
		Key:      "",
		Provider: "",
	}
	yamlFile, err := os.ReadFile(sk.File)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &cy)
	if err != nil {
		return nil, err
	}
	return &endefi.Secret{
		Key:      cy.Key,
		Provider: cy.Provider,
	}, nil
}

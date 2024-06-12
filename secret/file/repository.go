package file

import (
	"os"

	"github.com/bzhtux/endefi/config"
	"github.com/bzhtux/endefi/endefi"

	"gopkg.in/yaml.v2"
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
	sk := config.NewSecretConfig()
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

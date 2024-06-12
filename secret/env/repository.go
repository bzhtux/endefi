package env

import (
	"github.com/bzhtux/endefi/config"
	"github.com/bzhtux/endefi/endefi"
)

type EnvRepository struct {
	cfg *config.SecretConfig
}

func NewEnvRepository(cfg config.Config) endefi.SecretRepository {
	return &EnvRepository{
		cfg: cfg.Secret,
	}
}

func (r *EnvRepository) GetSecretKey(cfg *config.Config) (*endefi.Secret, error) {
	sk := config.NewSecretConfig()
	return &endefi.Secret{
		Key:      sk.Key,
		Provider: "env",
	}, nil
}

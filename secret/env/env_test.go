package env_test

import (
	"os"

	"github.com/bzhtux/endefi/config"
	"github.com/bzhtux/endefi/secret/env"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Env", func() {

	ginkgo.Describe("Test Env Provider", func() {
		ginkgo.Context("Test GetSecretKey", func() {
			ginkgo.It("with env provider", func() {
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_PROVIDER")
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_KEY")
				os.Setenv(config.ENV_PREFIX+"_SECRET_PROVIDER", "env")
				os.Setenv(config.ENV_PREFIX+"_SECRET_KEY", "my secret test key")
				cfg, err1 := config.NewConfig()
				gomega.Expect(err1).To(gomega.BeNil())
				repo := env.NewEnvRepository(*cfg)
				secret, err2 := repo.GetSecretKey(cfg)
				gomega.Expect(err2).To(gomega.BeNil())
				gomega.Expect(secret.Key).To(gomega.Equal("my secret test key"))
			})
		})
	})
})

package config_test

import (
	"os"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/bzhtux/endefi/config"
	"gopkg.in/yaml.v2"
)

type NeverExits struct{}

func (e NeverExits) ExitCode() int {
	return -1
}

type testSecretFile struct {
	key      string
	provider string
}

var _ = ginkgo.Describe("Config", func() {
	// var cfg *config.Config

	ginkgo.Describe("Test Config Provider", func() {
		ginkgo.Context("Without Env Var", func() {
			ginkgo.It("config.NewConfig should exit with err", func() {
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_PROVIDER")
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_KEY")
				cfg, err := config.NewConfig()
				gomega.Expect(err).NotTo(gomega.BeNil())
				gomega.Expect(cfg).To(gomega.BeNil())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("no secret provider found via environment variable"))
			})
		})
		ginkgo.Context("With Env Var", func() {
			ginkgo.It("config.NewConfig should not exit with err", func() {
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_PROVIDER")
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_KEY")
				os.Setenv(config.ENV_PREFIX+"_SECRET_PROVIDER", "env")
				os.Setenv(config.ENV_PREFIX+"_SECRET_KEY", "my secret test key")
				cfg, err := config.NewConfig()
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(cfg).NotTo(gomega.BeNil())
				gomega.Expect(cfg.App.Provider).To(gomega.Equal("env"))
			})
		})
		ginkgo.Context("With Unexpected Env Var (secret provider set to 'test')", func() {
			ginkgo.It("config.NewConfig should exit with err", func() {
				os.Setenv(config.ENV_PREFIX+"_SECRET_PROVIDER", "test")
				cfg, err := config.NewConfig()
				gomega.Expect(err).NotTo(gomega.BeNil())
				gomega.Expect(cfg).To(gomega.BeNil())
			})
		})
	})
	ginkgo.Describe("Test local provider", func() {
		ginkgo.Context("With env var but without an existing secret file", func() {
			ginkgo.It("It should return an error", func() {
				os.Setenv(config.ENV_PREFIX+"_SECRET_PROVIDER", "local")
				os.Setenv(config.ENV_PREFIX+"_SECRET_FILE", "/tmp/test.yaml")
				cfg, err := config.NewConfig()
				gomega.Expect(err).NotTo(gomega.BeNil())
				gomega.Expect(cfg).To(gomega.BeNil())
			})
		})
		ginkgo.Context("With env var and with an existing secret file", func() {
			ginkgo.It("It should not return an error", func() {
				d, err := os.MkdirTemp("", "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(d)
				f, err := os.CreateTemp(d, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer f.Close()
				os.Setenv(config.ENV_PREFIX+"_SECRET_PROVIDER", "local")
				os.Setenv(config.ENV_PREFIX+"_SECRET_FILE", f.Name())
				cfg, err := config.NewConfig()
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(cfg).NotTo(gomega.BeNil())
				gomega.Expect(cfg.Secret.File).To(gomega.Equal(f.Name()))
			})
		})
		ginkgo.Context("With bad content for secret file", func() {
			ginkgo.It("It should return an error", func() {
				d, err := os.MkdirTemp("", "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(d)
				f, err := os.CreateTemp(d, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer f.Close()
				_, err = f.WriteString("bad content")
				gomega.Expect(err).To(gomega.BeNil())

				os.Setenv(config.ENV_PREFIX+"_SECRET_PROVIDER", "local")
				os.Setenv(config.ENV_PREFIX+"_SECRET_FILE", f.Name())
				cfg, err := config.NewConfig()
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(cfg).NotTo(gomega.BeNil())
			})
		})
		ginkgo.Context("With good content for secret file", func() {
			ginkgo.It("It should NOT return an error", func() {
				d, err := os.MkdirTemp("", "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(d)
				f, err := os.CreateTemp(d, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer f.Close()
				enc := yaml.NewEncoder(f)
				err = enc.Encode(testSecretFile{
					key:      "testKey",
					provider: "local",
				})
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

})

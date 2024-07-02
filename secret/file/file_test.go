package file_test

import (
	"os"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"github.com/bzhtux/endefi/config"
	"github.com/bzhtux/endefi/secret/file"
)

type testSecretFile struct {
	Key      string
	Provider string
}

var _ = ginkgo.Describe("File", func() {

	ginkgo.Describe("Test File/Local Provider", func() {
		ginkgo.Context("Test GetSecretKey", func() {
			ginkgo.It("with local provider", func() {
				d, err := os.MkdirTemp("", "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(d)
				f, err := os.CreateTemp(d, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer f.Close()
				// _file, err := os.OpenFile(f.Name(), os.O_RDWR, 0644)
				// gomega.Expect(err).To(gomega.BeNil())
				// defer _file.Close()
				enc := yaml.NewEncoder(f)
				err = enc.Encode(testSecretFile{
					Key:      "test",
					Provider: "local",
				})
				gomega.Expect(err).To(gomega.BeNil())
				// fmt.Printf("CFG File: %s\n", f.Name())
				// out, err := exec.Command("cat", f.Name()).Output()
				// gomega.Expect(err).To(gomega.BeNil())
				// fmt.Printf("cat file: %s\n", string(out))
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_PROVIDER")
				os.Setenv(config.ENV_PREFIX+"_SECRET_PROVIDER", "local")
				os.Setenv(config.ENV_PREFIX+"_SECRET_FILE", f.Name())
				cfg, err1 := config.NewConfig()
				cfg.Secret.File = f.Name()
				gomega.Expect(err1).To(gomega.BeNil())
				repo := file.NewFileRepository(*cfg)
				secret, err2 := repo.GetSecretKey(cfg)
				gomega.Expect(err2).To(gomega.BeNil())
				gomega.Expect(secret.Key).To(gomega.Equal("test"))
			})
			ginkgo.It("with non existing file should return an error", func() {
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_PROVIDER")
				os.Setenv(config.ENV_PREFIX+"_SECRET_PROVIDER", "local")
				os.Setenv(config.ENV_PREFIX+"_SECRET_FILE", "/endefi-test")
				// _, err1 := config.NewConfig()
				// gomega.Expect(err1).NotTo(gomega.BeNil())
				ac := config.AppConfig{
					Provider: "local",
				}
				sc := config.SecretConfig{
					File: "/endefi-test",
				}
				// c := config.Config{
				// 	App:    &ac,
				// 	Secret: &sc,
				// }
				repo := file.NewFileRepository(config.Config{
					App:    &ac,
					Secret: &sc,
				})
				secret, err2 := repo.GetSecretKey(&config.Config{
					App:    &ac,
					Secret: &sc,
				})
				gomega.Expect(err2).NotTo(gomega.BeNil())
				gomega.Expect(secret).To(gomega.BeNil())
			})
		})
	})
})

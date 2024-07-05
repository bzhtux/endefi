package endefi_test

import (
	"os"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/bzhtux/endefi/config"
	"github.com/bzhtux/endefi/endefi"
)

type testSecretFile struct {
	key      string
	provider string
}

var _ = ginkgo.Describe("Endefi", func() {
	ginkgo.Describe("Test Endefi File", func() {
		ginkgo.Context("Test NewFile with a valid file", func() {
			ginkgo.It("should read a file without error", func() {
				testFile := endefi.File{}
				d, err := os.MkdirTemp("", "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(d)
				f, err := os.CreateTemp(d, "endefi-test")
				testFile.Path = f.Name()
				gomega.Expect(err).To(gomega.BeNil())
				defer f.Close()
				// _, err = f.WriteString("some content")
				// gomega.Expect(err).To(gomega.BeNil())
				_, err = endefi.NewFile(&testFile)
				gomega.Expect(err).To(gomega.BeNil())
			})
			ginkgo.It("should return file content without error", func() {
				testFile := endefi.File{}
				d, err := os.MkdirTemp("", "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(d)
				f, err := os.CreateTemp(d, "endefi-test")
				testFile.Path = f.Name()
				gomega.Expect(err).To(gomega.BeNil())
				defer f.Close()
				_, err = f.WriteString("some content")
				gomega.Expect(err).To(gomega.BeNil())
				ft, err := endefi.NewFile(&testFile)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(ft.Data).To(gomega.Equal([]byte("some content")))
			})
			ginkgo.Context("Test NewFile with an invalid file", func() {
				ginkgo.It("should return an error", func() {
					testFile := endefi.File{}
					testFile.Path = "invalid-path"
					_, err := endefi.NewFile(&testFile)
					gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
				})
			})
		})
	})
	ginkgo.Describe("Test Endefi FS", func() {
		ginkgo.Context("Test ListDir", func() {
			ginkgo.It("should list files without error", func() {
				d, err := os.MkdirTemp("", "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(d)
				f1, err := os.CreateTemp(d, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer f1.Close()
				f2, err := os.CreateTemp(d, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer f2.Close()
				_, err = f1.WriteString("some content")
				gomega.Expect(err).To(gomega.BeNil())
				_, err = f2.WriteString("some content")
				gomega.Expect(err).To(gomega.BeNil())
				files, err := endefi.ListDir(d)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(files).To(gomega.HaveLen(2))
				gomega.Expect(files).To(gomega.ContainElement(f1.Name()))
				gomega.Expect(files).To(gomega.ContainElement(f2.Name()))
			})
			ginkgo.It("with a fake dir should return an error", func() {
				_, err := endefi.ListDir("endefi-test")
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
		ginkgo.Context("Test WalkDir", func() {
			ginkgo.It("should list files without error", func() {
				d1, err := os.MkdirTemp("", "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(d1)
				f1, err := os.CreateTemp(d1, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer f1.Close()
				f2, err := os.CreateTemp(d1, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer f2.Close()
				_, err = f1.WriteString("some content")
				gomega.Expect(err).To(gomega.BeNil())
				_, err = f2.WriteString("some content")
				gomega.Expect(err).To(gomega.BeNil())

				d2, err := os.MkdirTemp(d1, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(d1)
				f3, err := os.CreateTemp(d2, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer f3.Close()
				f4, err := os.CreateTemp(d2, "endefi-test")
				gomega.Expect(err).To(gomega.BeNil())
				defer f4.Close()
				_, err = f3.WriteString("some content")
				gomega.Expect(err).To(gomega.BeNil())
				_, err = f4.WriteString("some content")
				gomega.Expect(err).To(gomega.BeNil())

				files, err := endefi.WalkDir(d1)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(files).To(gomega.HaveLen(4))
				gomega.Expect(files).To(gomega.ContainElement(f1.Name()))
				gomega.Expect(files).To(gomega.ContainElement(f2.Name()))
				gomega.Expect(files).To(gomega.ContainElement(f3.Name()))
				gomega.Expect(files).To(gomega.ContainElement(f4.Name()))
			})
			ginkgo.It("with a fake dir should return an error", func() {
				_, err := endefi.WalkDir("endefi-test")
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
	ginkgo.Describe("Test Endefi Logic", func() {
		ginkgo.Context("Test GenerateRandomKey", func() {
			ginkgo.It("should generate a random key without error", func() {
				_, err := endefi.GenerateRandomKey()
				gomega.Expect(err).To(gomega.BeNil())
			})
			ginkgo.It("should generate a 32 byte key without error", func() {
				new_key, _ := endefi.GenerateRandomKey()
				gomega.Expect(new_key).To(gomega.HaveLen(32))
			})
		})
		ginkgo.Context("Test key size", func() {
			ginkgo.It("with a generated random 32 byte key should not return an error", func() {
				new_key, _ := endefi.GenerateRandomKey()
				err := endefi.CheckKeySize([]byte(new_key))
				gomega.Expect(err).To(gomega.BeNil())
			})
			ginkgo.It("with a non 32 byte key should return an error", func() {
				err := endefi.CheckKeySize([]byte("test-key"))
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
			ginkgo.It("with a non 8 nor 16 nor 32 byte key should return an error", func() {
				err := endefi.CheckKeySize([]byte("mysecretkey"))
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
			ginkgo.It("with a 32 byte key should not return an error", func() {
				err := endefi.CheckKeySize([]byte("azertyuiopqsdfghjklmwxcvbn012345"))
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
		ginkgo.Context("Test Encrypt Data", func() {
			ginkgo.It("With a non valid key should return an error", func() {
				_, err := endefi.EncryptData("test", []byte("test"))
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
			ginkgo.It("With a valid key should not return an error", func() {
				new_key, err1 := endefi.GenerateRandomKey()
				gomega.Expect(err1).To(gomega.BeNil())
				encrypted, err2 := endefi.EncryptData("test", new_key)
				gomega.Expect(err2).To(gomega.BeNil())
				gomega.Expect(encrypted).NotTo(gomega.Equal("test"))
			})
		})
		ginkgo.Context("Test Decrypt Data", func() {
			ginkgo.It("With a non valid key should return an error", func() {
				_, err := endefi.DecryptData([]byte("test"), []byte("test"))
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
			ginkgo.It("With a valid key and plain data should return an error", func() {
				new_key, err1 := endefi.GenerateRandomKey()
				gomega.Expect(err1).To(gomega.BeNil())
				_, err2 := endefi.DecryptData([]byte("this is my encrypted test data"), new_key)
				gomega.Expect(err2).NotTo(gomega.BeNil())
			})
			ginkgo.It("With a valid key and encrypted data should not return an error", func() {
				new_key, err1 := endefi.GenerateRandomKey()
				gomega.Expect(err1).To(gomega.BeNil())
				encrypted, err2 := endefi.EncryptData("test plain data", new_key)
				gomega.Expect(err2).To(gomega.BeNil())
				plain, err3 := endefi.DecryptData(encrypted, new_key)
				gomega.Expect(err3).To(gomega.BeNil())
				gomega.Expect(plain).To(gomega.Equal([]byte("test plain data")))
			})
			ginkgo.It("With a valid key but not the one used to encrpypt data should return an error", func() {
				new_key1, err1 := endefi.GenerateRandomKey()
				gomega.Expect(err1).To(gomega.BeNil())
				encrypted, err2 := endefi.EncryptData("test plain data", new_key1)
				gomega.Expect(err2).To(gomega.BeNil())
				new_key2, err3 := endefi.GenerateRandomKey()
				gomega.Expect(err3).To(gomega.BeNil())
				_, err4 := endefi.DecryptData(encrypted, new_key2)
				gomega.Expect(err4).NotTo(gomega.BeNil())
			})
		})
	})
	ginkgo.Describe("Test Endefi Secret", func() {
		ginkgo.Context("Test Encrypt Secret file", func() {
			ginkgo.It("With non existent file should return an error", func() {
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_PROVIDER")
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_FILE")
				os.Setenv(config.ENV_PREFIX+"_SECRET_PROVIDER", "env")
				os.Setenv(config.ENV_PREFIX+"_SECRET_FILE", "/tmp/.fake/secret.fake")
				err := endefi.EncryptSecretFile(config.ENV_PREFIX+"_SECRET_FILE", []byte("test"))
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
			ginkgo.It("With existing file should not return an error", func() {
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_PROVIDER")
				os.Unsetenv(config.ENV_PREFIX + "_SECRET_FILE")
				os.Setenv(config.ENV_PREFIX+"_SECRET_PROVIDER", "env")
				d, err := os.MkdirTemp(os.TempDir(), ".fake")
				gomega.Expect(err).To(gomega.BeNil())
				defer os.RemoveAll(d)
				f, err := os.CreateTemp(d, "secret.fake")
				// fmt.Printf("File name: %s\n", f.Name())
				gomega.Expect(err).To(gomega.BeNil())
				os.Setenv(config.ENV_PREFIX+"_SECRET_FILE", f.Name())
				err = endefi.EncryptSecretFile(f.Name(), []byte("test"))
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})

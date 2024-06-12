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
package cmd

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/bzhtux/endefi/endefi"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a local file",
	Long:  `Decrypt a local file using aes GCM.`,
	Run: func(cmd *cobra.Command, args []string) {
		if dirPath != "" {
			if recursive {
				log.Default().Printf("Decrypt all files recursively within %s", dirPath)
				files, err := endefi.WalkDir(dirPath)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
				for _, filePath := range files {
					log.Default().Printf("Decrypt a new file: %s", filePath)
					if err := DecryptFile(filePath); err != nil {
						log.Fatal(err)
						os.Exit(1)
					}
				}
			} else {
				log.Default().Printf("Decrypt a new directory: %s", dirPath)
				files, err := endefi.ListDir(dirPath)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
				for _, filePath := range files {
					if err := DecryptFile(filePath); err != nil {
						log.Fatal(err)
						os.Exit(1)
					}
				}
			}

		} else {
			log.Default().Printf("Decrypt a new file: %s", filePath)
			if err := DecryptFile(filePath); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		}
	},
}

// var filePath string

func init() {
	rootCmd.AddCommand(decryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")
	decryptCmd.Flags().StringVarP(&filePath, "file", "f", "", "Local file to decrypt")
	decryptCmd.Flags().StringVarP(&dirPath, "dir", "d", "", "Local directory to encrypt")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func DecryptFile(filePath string) error {
	info, err := os.Lstat(filePath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	f := &endefi.File{
		Path: filePath,
	}
	f, err = endefi.NewFile(f)
	if err != nil {
		log.Fatal(err)
		return err
	}
	k, err := service.GetSecretKey(cfg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	key, err := hex.DecodeString(k.Key)
	if err != nil {
		log.Fatal(err)
		return err
	}
	decryptedData, err := endefi.DecryptData(f.Data, key)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if err := os.WriteFile(filePath, decryptedData, info.Mode().Perm()); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

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
	"log"
	"os"

	"github.com/bzhtux/endefi/config"
	"github.com/bzhtux/endefi/endefi"
	"github.com/bzhtux/endefi/secret/env"
	"github.com/bzhtux/endefi/secret/file"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Endefi",
	Short: "Endefi stands for Encrypt Decrypt Files",
	Long: `Endefi is a command line tool to encrypt and decrypt local files using 
multiple external secret operators to get the secret key required to encrypt and decrypt files.
	
Run 'endefi --help for more informations
	
Get more details at github.com/bzhtux/endefi'
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var filePath string
var cfg *config.Config
var service endefi.SecretService

// var repo endefi.SecretRepository

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cfg = config.NewConfig()
	// fmt.Printf("cfg: %v\n", cfg.App.Provider)
	repo := selectSecretRepo(cfg)
	service = endefi.NewSecretService(repo)

	log.Default().Printf("Starting %s", cfg.App.Name)
	// cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.endefi.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func selectSecretRepo(cfg *config.Config) endefi.SecretRepository {
	switch cfg.App.Provider {
	case "env":
		// repo = env.NewEnvRepository(cfg.Secret)
		repo := env.NewEnvRepository(*cfg)
		return repo
	case "local":
		// repo := bw.NewBitwardenRepository(*cfg)
		repo := file.NewFileRepository(*cfg)
		return repo
	default:
		log.Fatal("Unknown secret provider")
		os.Exit(1)
	}

	return nil
}

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"log"
	"os"
	"path/filepath"

	"github.com/bzhtux/endefi/config"
	"github.com/bzhtux/endefi/endefi"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type EndefiSecretFile struct {
	Key      string
	Provider string
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Create $HOME/.endefi/secret.yaml with random secret key",
	Long:  `Create $HOME/.endefi/secret.yaml with random secret key`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			log.Fatal(err)
		}
		if cfg.Secret.File != "" {
			log.Default().Printf("Secret file: %s", cfg.Secret.File)
		}
		// homedir, err := homedir.Dir()
		// if err != nil {
		// 	log.Fatal("Homedir error:", err.Error())
		// }
		_, err = os.Stat(cfg.Secret.File)
		if err != nil {
			if err := os.Mkdir(filepath.Dir(cfg.Secret.File), 0700); err != nil {
				log.Fatal("Mkdir error:", err.Error())
			}
			secretkey, err := endefi.GenerateRandomKey()
			if err != nil {
				log.Fatal(err.Error())
			}
			esf := EndefiSecretFile{
				Key:      hex.EncodeToString(secretkey),
				Provider: "local",
			}

			f, err := os.Create(cfg.Secret.File)
			if err != nil {
				log.Fatal("Create error:", err.Error())
				log.Fatal(err)
			}
			enc := yaml.NewEncoder(f)
			if err := enc.Encode(esf); err != nil {
				log.Fatal("Encode error:", err.Error())
			} else {
				log.Printf("Secret file created: %s", cfg.Secret.File)
			}
		} else {
			log.Printf("Secret file already exists: %s", cfg.Secret.File)
		}

	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package cmd

import (
	"log"
	"strings"

	"github.com/SawitProRecruitment/UserService/lib/locker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(encryptCmd)
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Args:  cobra.ExactArgs(1),
	Short: "Encrypt text with AES to help create valid password",
	Run: func(cmd *cobra.Command, args []string) {
		plainText := args[0]
		trimmed := strings.TrimSpace(plainText)

		cfg := initConfig()
		keyLock := locker.New(cfg.AES.SecretKey)
		result, err := keyLock.Encrypt(trimmed)
		if err != nil {
			log.Fatalf("err: %v", err)
		}

		log.Printf("plain text: %s\n", plainText)
		log.Println("encrypted result:")
		log.Println(result)
	},
}

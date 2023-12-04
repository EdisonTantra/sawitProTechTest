package cmd

import (
	"fmt"
	"strings"

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
		secretKey := cfg.Auth.EncryptSecretKey
		res := encryptText(trimmed, secretKey)

		fmt.Printf("plain text: %s\n", plainText)
		fmt.Println("result:")
		fmt.Println(res)
	},
}

func encryptText(text string, key string) string {
	//TODO implement
	return "result"
}

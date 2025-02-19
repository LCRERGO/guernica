package cmd

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/LCRERGO/guernica/pkg/alphabet"
	"github.com/LCRERGO/guernica/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "guernica",
		Short: "A password generator",
		Long:  `Guernica is a password generator that provides ways of setting the alphabet and the password length.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("guernica %s\n", "v1.0.0")
			alphabet := alphabet.GetAlphabet(config.GetConfig().PasswordAlphabet)
			genPassword := generatePassword(config.GetConfig().Length, alphabet)
			fmt.Printf("Password: %s", genPassword)

		},
	}
)

func init() {
	rootCmd.PersistentFlags().IntVarP(&config.GetConfig().Length,
		"length", "l", 20, "password length")
	rootCmd.PersistentFlags().StringVarP(&config.GetConfig().PasswordAlphabet,
		"alphabet", "a", "full",
		"password alphabet [digits,lower,upper,alnum,full,complete]")
	viper.SetDefault("author", "Lucas Cruz dos Reis <lcr.ergo@gmail.com>")
	viper.SetDefault("license", "MIT")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generatePassword(length int, alphabet string) string {
	password := ""
	println(alphabet)
	alphabetLength := len(alphabet)
	runedAlphabet := []rune(alphabet)

	for i := 0; i < length; i++ {
		randInt, err := rand.Int(rand.Reader, big.NewInt(int64(alphabetLength)))
		if err != nil {
			log.Print(err)
		}
		password += string(runedAlphabet[randInt.Int64()])
	}

	return password
}

package cmd

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/LCRERGO/guernica/pkg/alphabet"
	"github.com/LCRERGO/guernica/pkg/config"
	"github.com/LCRERGO/guernica/pkg/passphrase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "v1.1.0"

var rootCmd = &cobra.Command{
	Use:   "guernica",
	Short: "A password generator",
	Long:  `Guernica is a password generator that provides ways of setting the alphabet and the password length.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.GetConfig()

		if err := validate(cmd, cfg); err != nil {
			return err
		}

		fmt.Printf("guernica %s\n", version)

		if cfg.Words > 0 {
			return runPassphrases(cfg)
		}

		return runPasswords(cfg)
	},
}

func init() {
	cfg := config.GetConfig()
	rootCmd.PersistentFlags().IntVarP(&cfg.Length,
		"length", "l", 20, "password length")
	rootCmd.PersistentFlags().StringVarP(&cfg.PasswordAlphabet,
		"alphabet", "a", "full",
		"password alphabet [digits,lower,upper,alnum,full,complete]")
	rootCmd.PersistentFlags().StringVar(&cfg.CustomAlphabet,
		"custom-alphabet", "", "use an arbitrary character set as alphabet")
	rootCmd.PersistentFlags().IntVarP(&cfg.Count,
		"count", "n", 1, "number of passwords to generate")
	rootCmd.PersistentFlags().BoolVar(&cfg.Loop,
		"loop", false, "generate passwords forever until interrupted")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Clipboard,
		"clipboard", "c", false, "copy the generated password to the clipboard")
	rootCmd.PersistentFlags().IntVarP(&cfg.Words,
		"words", "w", 0, "generate a passphrase of N words instead of a password")
	rootCmd.PersistentFlags().StringVarP(&cfg.Separator,
		"separator", "s", " ", "word separator for passphrases")
	rootCmd.PersistentFlags().StringVar(&cfg.WordlistFile,
		"wordlist-file", "", "path to a custom wordlist for passphrases")
	rootCmd.PersistentFlags().BoolVarP(&cfg.NoEntropy,
		"no-entropy", "e", false, "do not print the entropy of generated passwords")
	rootCmd.PersistentFlags().BoolVarP(&cfg.RequireAll,
		"require-all", "r", false, "ensure at least one character from every alphabet class")
	rootCmd.PersistentFlags().BoolVarP(&cfg.NoAmbiguous,
		"no-ambiguous", "x", false, "exclude ambiguous characters (l, 1, I, O, 0, ...)")
	rootCmd.PersistentFlags().StringVar(&cfg.AmbiguousExclude,
		"ambiguous-exclude", "", "custom character set to strip when excluding ambiguous characters")
	viper.SetDefault("author", "Lucas Cruz dos Reis <lcr.ergo@gmail.com>")
	viper.SetDefault("license", "MIT")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// validate rejects conflicting flag combinations before any password is
// generated.
func validate(cmd *cobra.Command, cfg *config.Config) error {
	if cfg.Count < 1 {
		return fmt.Errorf("count must be at least 1")
	}
	if cfg.Length < 1 {
		return fmt.Errorf("length must be at least 1")
	}

	if cfg.Loop && cfg.Count != 1 {
		return fmt.Errorf("--loop and --count cannot be combined")
	}
	if cfg.Clipboard && cfg.Loop {
		return fmt.Errorf("--clipboard and --loop cannot be combined")
	}
	if cfg.Clipboard && cfg.Count != 1 {
		return fmt.Errorf("--clipboard and --count cannot be combined")
	}

	if cfg.CustomAlphabet != "" && cmd.Flags().Changed("alphabet") {
		return fmt.Errorf("--custom-alphabet and --alphabet cannot be combined")
	}
	if cfg.RequireAll && cfg.CustomAlphabet != "" {
		return fmt.Errorf("--require-all cannot be used with --custom-alphabet")
	}

	if cfg.RequireAll {
		classes, _ := alphabet.GetAlphabetClasses(cfg.PasswordAlphabet)
		if len(classes) > cfg.Length {
			return fmt.Errorf("length %d is too small to include all %d alphabet classes", cfg.Length, len(classes))
		}
	}

	if cfg.Words > 0 {
		if cmd.Flags().Changed("length") {
			return fmt.Errorf("--words and --length cannot be combined")
		}
		if cmd.Flags().Changed("alphabet") || cfg.CustomAlphabet != "" {
			return fmt.Errorf("--words cannot be combined with alphabet flags")
		}
		if cfg.RequireAll || cfg.NoAmbiguous {
			return fmt.Errorf("--words cannot be combined with --require-all or --no-ambiguous")
		}
	}

	return nil
}

func runPasswords(cfg *config.Config) error {
	var alphabetStr string
	var classes []string

	if cfg.CustomAlphabet != "" {
		alphabetStr = cfg.CustomAlphabet
	} else {
		alphabetStr = alphabet.GetAlphabet(cfg.PasswordAlphabet)
		classes, _ = alphabet.GetAlphabetClasses(cfg.PasswordAlphabet)
	}

	if cfg.NoAmbiguous || cfg.AmbiguousExclude != "" {
		exclusions := cfg.AmbiguousExclude
		if exclusions == "" {
			exclusions = alphabet.DefaultAmbiguous
		}
		alphabetStr = alphabet.Exclude(alphabetStr, exclusions)

		var filtered []string
		for _, cls := range classes {
			if c := alphabet.Exclude(cls, exclusions); c != "" {
				filtered = append(filtered, c)
			}
		}
		classes = filtered
	}

	printHeader := func(i int) string {
		if cfg.Loop || cfg.Count > 1 {
			return fmt.Sprintf("%d. ", i)
		}
		return "Password: "
	}

	for i := 1; cfg.Loop || i <= cfg.Count; i++ {
		pw, err := generatePassword(cfg.Length, alphabetStr, classes, cfg.RequireAll)
		if err != nil {
			return err
		}

		fmt.Printf("%s%s\n", printHeader(i), pw)
		if !cfg.NoEntropy {
			fmt.Printf("   Entropy: %.2f bits\n", entropyOf(len([]rune(alphabetStr)), cfg.Length))
		}

		if cfg.Clipboard {
			if err := copyToClipboard(pw); err != nil {
				return err
			}
		}
	}

	return nil
}

func runPassphrases(cfg *config.Config) error {
	words, err := passphrase.GetWordlist(cfg.WordlistFile)
	if err != nil {
		return err
	}
	if cfg.Words > len(words) {
		return fmt.Errorf("word count %d exceeds wordlist size %d", cfg.Words, len(words))
	}

	for i := 1; cfg.Loop || i <= cfg.Count; i++ {
		phrase, err := passphrase.Generate(cfg.Words, cfg.Separator, words)
		if err != nil {
			return err
		}

		if cfg.Loop || cfg.Count > 1 {
			fmt.Printf("%d. %s\n", i, phrase)
		} else {
			fmt.Printf("Passphrase: %s\n", phrase)
		}
		if !cfg.NoEntropy {
			fmt.Printf("   Entropy: %.2f bits\n", passphrase.Entropy(cfg.Words, len(words)))
		}

		if cfg.Clipboard {
			if err := copyToClipboard(phrase); err != nil {
				return err
			}
		}
	}

	return nil
}

func generatePassword(length int, alphabetStr string, classes []string, requireAll bool) (string, error) {
	runed := []rune(alphabetStr)
	password := make([]rune, 0, length)

	if requireAll {
		// Guarantee one character from every enabled class, then fill the rest
		// uniformly and shuffle.
		for _, cls := range classes {
			r, err := randomRune([]rune(cls))
			if err != nil {
				return "", err
			}
			password = append(password, r)
		}
	}

	for len(password) < length {
		r, err := randomRune(runed)
		if err != nil {
			return "", err
		}
		password = append(password, r)
	}

	if requireAll {
		shuffle(password)
	}

	return string(password), nil
}

func randomRune(set []rune) (rune, error) {
	n := int64(len(set))
	randInt, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		return 0, err
	}
	return set[randInt.Int64()], nil
}

func shuffle(runes []rune) {
	for i := len(runes) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			continue
		}
		runes[i], runes[j.Int64()] = runes[j.Int64()], runes[i]
	}
}

func entropyOf(alphabetSize, length int) float64 {
	return float64(length) * math.Log2(float64(alphabetSize))
}

// copyToClipboard pipes text into the platform clipboard tool. It returns an
// error when no tool is available.
func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	default:
		switch {
		case lookup("wl-copy"):
			cmd = exec.Command("wl-copy")
		case lookup("xclip"):
			cmd = exec.Command("xclip", "-selection", "clipboard")
		case lookup("xsel"):
			cmd = exec.Command("xsel", "--clipboard", "--input")
		default:
			return fmt.Errorf("no clipboard tool available (install wl-copy, xclip, or xsel)")
		}
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func lookup(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

package passphrase

import (
	"bufio"
	"crypto/rand"
	_ "embed"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
)

//go:embed eff_large_wordlist.txt
var defaultWordlist string

// GetWordlist returns the words from the embedded EFF large wordlist, or from
// the file at the given path when one is provided.
func GetWordlist(wordlistFile string) ([]string, error) {
	if wordlistFile == "" {
		return strings.Fields(defaultWordlist), nil
	}

	f, err := os.Open(wordlistFile)
	if err != nil {
		return nil, fmt.Errorf("open wordlist: %w", err)
	}
	defer f.Close()

	var words []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if word := strings.TrimSpace(scanner.Text()); word != "" {
			words = append(words, word)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read wordlist: %w", err)
	}

	return words, nil
}

// Generate produces a passphrase of n words joined by separator, picking each
// word uniformly from words using a cryptographically secure RNG.
func Generate(n int, separator string, words []string) (string, error) {
	if len(words) == 0 {
		return "", fmt.Errorf("empty wordlist")
	}

	selected := make([]string, n)
	max := big.NewInt(int64(len(words)))
	for i := 0; i < n; i++ {
		randInt, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("random: %w", err)
		}
		selected[i] = words[randInt.Int64()]
	}

	return strings.Join(selected, separator), nil
}

// Entropy returns the theoretical entropy in bits of a passphrase drawn from a
// wordlist of size listSize with n words.
func Entropy(n, listSize int) float64 {
	return float64(n) * math.Log2(float64(listSize))
}

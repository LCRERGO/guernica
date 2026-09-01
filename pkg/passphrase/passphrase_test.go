package passphrase

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWordlistEmbedded(t *testing.T) {
	words, err := GetWordlist("")
	require.NoError(t, err)
	assert.Equal(t, 7776, len(words))
	assert.Contains(t, words, "abacus")
	assert.Contains(t, words, "zoom")
}

func TestGetWordlistFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "words.txt")
	require.NoError(t, os.WriteFile(path, []byte("alpha\nbeta\n\ngamma\n"), 0o644))

	words, err := GetWordlist(path)
	require.NoError(t, err)
	assert.Equal(t, []string{"alpha", "beta", "gamma"}, words)
}

func TestGetWordlistMissingFile(t *testing.T) {
	_, err := GetWordlist("/nonexistent/path/words.txt")
	assert.Error(t, err)
}

func TestGenerate(t *testing.T) {
	words := []string{"alpha", "beta", "gamma"}

	phrase, err := Generate(4, "-", words)
	require.NoError(t, err)

	parts := strings.Split(phrase, "-")
	assert.Len(t, parts, 4)
	for _, p := range parts {
		assert.Contains(t, words, p)
	}
}

func TestGenerateEmptyWordlist(t *testing.T) {
	_, err := Generate(3, " ", nil)
	assert.Error(t, err)
}

func TestEntropy(t *testing.T) {
	assert.InDelta(t, 12.9248, Entropy(1, 7776), 0.001)
	assert.InDelta(t, 64.6241, Entropy(5, 7776), 0.001)
	assert.InDelta(t, 0.0, Entropy(0, 7776), 0.001)
}

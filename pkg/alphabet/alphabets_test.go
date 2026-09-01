package alphabet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildAlphabets(t *testing.T) {
	tests := []struct {
		alphaID int
		want    string
	}{
		{int(Digits), "0123456789"},
		{int(Digits | LowercaseLetters), "0123456789abcdefghijklmnopqrstuvwxyz"},
		{int(Digits | UppercaseLetters), "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{int(Alnum), "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{int(Full), "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!%@*+$^&#"},
		{int(Complete), "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!%@*+$^&#{?)=\\\"]>(_-<;}~`\\':.,[|/"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, buildAlphabet(tt.alphaID))
	}
}

func TestExclude(t *testing.T) {
	tests := []struct {
		alphabet   string
		exclusions string
		want       string
	}{
		{"0123456789", "", "0123456789"},
		{"0123456789", "05", "12346789"},
		{"abcdefghijklmnopqrstuvwxyz", "aeiou", "bcdfghjklmnpqrstvwxyz"},
		{"ABCabc", "ABC", "abc"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, Exclude(tt.alphabet, tt.exclusions))
	}
}

func TestGetAlphabetClasses(t *testing.T) {
	tests := []struct {
		preset string
		want   []string
	}{
		{"digits", []string{"0123456789"}},
		{"alnum", []string{"0123456789", "abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"}},
		{"full", []string{"0123456789", "abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "!%@*+$^&#"}},
	}

	for _, tt := range tests {
		classes, err := GetAlphabetClasses(tt.preset)
		assert.NoError(t, err)
		assert.Equal(t, tt.want, classes)
	}

	_, err := GetAlphabetClasses("bogus")
	assert.Error(t, err)
}

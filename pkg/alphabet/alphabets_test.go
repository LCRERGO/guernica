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

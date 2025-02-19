package alphabet

import (
	"fmt"
	"os"
	"strings"
)

type AlphabetType int

const alphabetMaxNum = 16

const (
	Digits AlphabetType = 0x1 << iota
	LowercaseLetters
	UppercaseLetters
	CommonSymbols
	AdditionalSymbols

	Alnum    = Digits | LowercaseLetters | UppercaseLetters
	Full     = Alnum | CommonSymbols
	Complete = Full | AdditionalSymbols
)

var alphabets = map[AlphabetType]string{
	Digits:            "0123456789",
	LowercaseLetters:  "abcdefghijklmnopqrstuvwxyz",
	UppercaseLetters:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	CommonSymbols:     "!%@*+$^&#",
	AdditionalSymbols: "{?)=\\\"]>(_-<;}~`\\':.,[|/",
}

func (t AlphabetType) String() string {
	return alphabets[t]
}

func GetAlphabet(alphaSymbol string) string {
	var alphabet string

	symbol := strings.Trim(strings.ToLower(alphaSymbol), " ")
	switch symbol {
	case "digits":
		alphabet = buildAlphabet(int(Digits))
	case "lower":
		alphabet = buildAlphabet(int(LowercaseLetters))
	case "upper":
		alphabet = buildAlphabet(int(UppercaseLetters))
	case "alnum":
		alphabet = buildAlphabet(int(Alnum))
	case "full":
		alphabet = buildAlphabet(int(Full))
	case "complete":
		alphabet = buildAlphabet(int(Complete))
	default:
		fmt.Printf("invalid alphabet: %v\n", alphaSymbol)
		os.Exit(1)
	}

	return alphabet
}

func buildAlphabet(alphaID int) string {
	var alphabet string

	for i := 0; i < alphabetMaxNum; i++ {
		if ((alphaID >> i) & 0x01) == 0x1 {
			alphabet += alphabets[0x1<<i]
		}
	}

	return alphabet
}

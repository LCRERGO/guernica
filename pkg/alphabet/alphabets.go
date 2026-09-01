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

// DefaultAmbiguous holds the characters that look alike and are commonly
// confused, excluded by --no-ambiguous unless overridden.
const DefaultAmbiguous = "0O1lI|`'\""

// Exclude strips every rune in exclusions from the given alphabet and returns
// the filtered string.
func Exclude(alphabet, exclusions string) string {
	if exclusions == "" {
		return alphabet
	}

	drop := make(map[rune]struct{}, len(exclusions))
	for _, r := range exclusions {
		drop[r] = struct{}{}
	}

	var b strings.Builder
	for _, r := range alphabet {
		if _, ok := drop[r]; !ok {
			b.WriteRune(r)
		}
	}

	return b.String()
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

// GetAlphabetClasses returns the individual class strings that make up the
// named preset. It is used by --require-all to guarantee at least one
// character from every enabled class.
func GetAlphabetClasses(alphaSymbol string) ([]string, error) {
	symbol := strings.Trim(strings.ToLower(alphaSymbol), " ")
	switch symbol {
	case "digits":
		return []string{alphabets[Digits]}, nil
	case "lower":
		return []string{alphabets[LowercaseLetters]}, nil
	case "upper":
		return []string{alphabets[UppercaseLetters]}, nil
	case "alnum":
		return []string{alphabets[Digits], alphabets[LowercaseLetters], alphabets[UppercaseLetters]}, nil
	case "full":
		return []string{alphabets[Digits], alphabets[LowercaseLetters], alphabets[UppercaseLetters], alphabets[CommonSymbols]}, nil
	case "complete":
		return []string{alphabets[Digits], alphabets[LowercaseLetters], alphabets[UppercaseLetters], alphabets[CommonSymbols], alphabets[AdditionalSymbols]}, nil
	default:
		return nil, fmt.Errorf("invalid alphabet: %v", alphaSymbol)
	}
}

func buildAlphabet(alphaID int) string {
	var alphabetBuilder strings.Builder

	for i := range alphabetMaxNum {
		if ((alphaID >> i) & 0x01) == 0x1 {
			alphabetBuilder.WriteString(alphabets[0x1<<i])
		}
	}

	return alphabetBuilder.String()
}

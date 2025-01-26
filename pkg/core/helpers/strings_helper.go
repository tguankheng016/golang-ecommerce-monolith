package helpers

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dustin/go-humanize/english"
)

// ExtractTextBetweenToPatterns extracts the text between fromPattern and toPattern from str.
//
// The function first searches for toPattern in str. If it is not found, an empty string is returned.
// After that, the function searches for the last occurrence of fromPattern in the substring before toPattern.
// If fromPattern is not found, the substring from the beginning to toPattern is returned. Otherwise, the
// string between fromPattern and toPattern is returned.
func ExtractTextBetweenToPatterns(str, fromPattern, toPattern string) string {
	num := strings.Index(str, toPattern)
	if num == -1 {
		return ""
	}

	text := str[:num]
	num2 := strings.LastIndex(text, fromPattern)
	if num2 == -1 {
		num2 = 0
		fromPattern = ""
	}

	return text[num2+len(fromPattern):]
}

// ReplaceLast replaces the last occurrence of pattern in str with replacement.
// If pattern is not found in str, str is returned unchanged.
func ReplaceLast(str, pattern, replacement string) string {
	if pattern == "" {
		return str
	}
	lastIndex := strings.LastIndex(str, pattern)
	if lastIndex == -1 {
		return str
	}
	return str[:lastIndex] + replacement + str[lastIndex+len(pattern):]
}

// Pluralize returns the plural form of a given string using the `english` package.
// It assumes that the string represents a singular noun and returns its plural form.
func Pluralize(str string) string {
	return english.PluralWord(2, str, "")
}

// ToUpperCaseFirstChar converts the first character of a string to uppercase.
// If the string is empty, it returns the string unchanged. The function checks
// if the first character is a letter before converting it to uppercase.
func ToUpperCaseFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	if unicode.IsLetter(r) {
		return string(unicode.ToUpper(r)) + s[n:]
	}
	return s
}

// Temporary convert the 1st new line to avoid replace
func ConvertDoubleNewLine(s string) string {
	return strings.Replace(s, "\n\n", "ConvertedNewLine\n", -1)
}

// Revert back the 1st new line back to \n
func RevertDoubleNewLine(s string) string {
	return strings.Replace(s, "ConvertedNewLine\n", "\n\n", -1)
}

func SplitCamelCase(str, splitter string) string {
	if !utf8.ValidString(str) {
		return str
	}
	entries := []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range str {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}

	return strings.Join(entries, splitter)
}

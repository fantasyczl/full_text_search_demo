package fts

import (
	"strings"
	"unicode"
)

func Tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func Analyze(text string) []string {
	tokens := Tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopWordsFilter(tokens)
	tokens = stemmerFilter(tokens)
	return tokens
}

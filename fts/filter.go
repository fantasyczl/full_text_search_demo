package fts

import "strings"
import snowballeng "github.com/kljensen/snowball/english"

type Filter func([]string) []string

func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))

	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}

	return r
}

var stopwords = map[string]bool{
	"a": true, "and": true, "be": true, "have": true, "i": true,
	"in": true, "of": true, "that": true, "the": true, "to": true,
}

func stopWordsFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if !stopwords[token] {
			r = append(r, token)
		}
	}

	return r
}

func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}

	return r
}

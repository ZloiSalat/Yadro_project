package pkg

import (
	"YadroProject/types"
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
)

type Stem interface {
	stemmedWords([]string) []string
}

type Stemmer struct {
}

func NewStemmer() *Stemmer {
	return &Stemmer{}
}

func stemWord(s string) string {
	return english.Stem(s, true)
}

func isNormalized(s string) bool {
	return !english.IsStopWord(s) && !strings.Contains(s, "'")

}

func (stm *Stemmer) stemmedWords(s types.Num) []string {
	m := make(map[string]bool)
	var word string
	var finalWords []string
	for i := 0; i < len(s.Keywords); i++ {
		word = stemWord(s.Keywords[i])
		if !isNormalized(word) || m[word] {
			continue
		}
		m[word] = true
		finalWords = append(finalWords, word)

	}

	return finalWords
}

func stringSplitter(s rune) bool {
	if !unicode.IsLetter(s) {
		return true
	}
	return false
}

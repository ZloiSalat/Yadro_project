package main

import (
	"flag"
	"fmt"
	"github.com/kljensen/snowball/english"
	"strings"
)

func stemWord(s string) string {
	return english.Stem(s, true)
}

func isNormalized(s string) bool {
	if !english.IsStopWord(s) && !strings.Contains(s, "'") {
		return true
	}
	return false
}

func stemmedWords(s []string) []string {
	var m = make(map[string]bool)
	var word string
	var finalWords []string
	for i := 0; i < len(s); i++ {
		word = stemWord(s[i])
		if !isNormalized(word) || m[word] {
			continue
		}
		m[word] = true
		finalWords = append(finalWords, word)

	}

	return finalWords
}

func stringSplitter(s rune) bool {
	switch s {
	case ' ', ':', ',', ';', '!', '?':
		return true
	}
	return false
}

func ReadCliCommands() []string {
	var sentence string
	flag.StringVar(&sentence, "s", "", "Write down the sentence")
	flag.Parse()

	return strings.FieldsFunc(sentence, stringSplitter)
}

func main() {
	words := ReadCliCommands()
	for _, word := range stemmedWords(words) {
		fmt.Println(word)
	}

}

package nlp

import (
	"iter"
	"regexp"
	"sort"
	"strings"
)

type WordCount map[string]int

var (
	wordRegex    = regexp.MustCompile(`\p{L}+`)
	prepositions = map[string]struct{}{
		"в": {}, "на": {}, "и": {}, "не": {}, "что": {}, "с": {}, "по": {}, "о": {}, "за": {}, "от": {},
		"у": {}, "к": {}, "до": {}, "из": {}, "без": {}, "для": {}, "при": {}, "про": {}, "об": {},
		"под": {}, "а": {}, "то": {}, "но": {}, "ли": {}, "же": {}, "ни": {}, "ну": {},
	}
	CyrillicCheck = regexp.MustCompile(`^[а-яА-ЯёЁ]+$`)
)

func Tokenize(line string) []string {
	tokens := wordRegex.FindAllString(line, -1)
	var result []string
	for _, token := range tokens {
		token = strings.ToLower(token)
		if _, stop := prepositions[token]; stop || !isRussian(token) {
			continue
		}
		result = append(result, normalize(token))
	}
	return result
}

func isRussian(word string) bool {
	return CyrillicCheck.MatchString(word)
}

func normalize(word string) string {
	return word
}

func Top(counts WordCount, topN int) iter.Seq2[string, int] {
	type kv struct {
		Word  string
		Count int
	}

	var sorted []kv
	for word, count := range counts {
		sorted = append(sorted, kv{word, count})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	return func(yield func(string, int) bool) {
		for _, pair := range sorted[:topN] {
			if !yield(pair.Word, pair.Count) {
				return
			}
		}
	}
}

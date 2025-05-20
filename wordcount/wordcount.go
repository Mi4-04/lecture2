package wordcount

import (
	"fmt"
	"sort"
)

type WordCount map[string]int

func PrintTopWords(counts WordCount, top int) {
	type kv struct {
		Word  string
		Count int
	}

	var sorted []kv
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	fmt.Printf("\nТоп %d слов:\n", top)
	for i, pair := range sorted[:top] {
		fmt.Printf("%2d. %s — %d\n", i+1, pair.Word, pair.Count)
	}
}

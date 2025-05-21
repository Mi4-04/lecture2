package wordcount

import (
	"fmt"
	"lecture2/nlp"
)

func PrintTopWords(counts nlp.WordCount, top int) {
	i := 1
	for word, count := range nlp.Top(counts, top) {
		fmt.Printf("%2d. %s — %d\n", i, word, count)
		i++
	}
}

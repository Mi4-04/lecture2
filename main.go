package main

import (
	"lecture2/utils"
	"lecture2/wordcount"
	"sync"
)

type Pair struct {
	Word  string
	Count int
}

const workerCount = 4
const topWordsCount = 100

func main() {
	files := []string{
		"tom-1.txt",
		"tom-2.txt",
		"tom-3.txt",
		"tom-4.txt",
	}

	lines := make(chan string, 1000)
	pairs := make(chan Pair, 1000)

	var readerWg sync.WaitGroup
	for _, file := range files {
		readerWg.Add(1)
		go func(f string) {
			defer readerWg.Done()
			utils.ReadFileLines(f, lines)
		}(file)
	}

	go func() {
		readerWg.Wait()
		close(lines)
	}()

	var workerWg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			mapper(lines, pairs)
		}()
	}

	go func() {
		workerWg.Wait()
		close(pairs)
	}()

	finalCounts := reducer(pairs)

	wordcount.PrintTopWords(finalCounts, topWordsCount)
}

func mapper(in <-chan string, out chan<- Pair) {
	for line := range in {
		words := utils.Tokenize(line)
		counted := make(map[string]int)

		for _, word := range words {
			if _, ok := utils.Prepositions[word]; ok || !utils.IsRussian(word) {
				continue
			}
			word = utils.Normalize(word)
			counted[word]++
		}

		for word, count := range counted {
			out <- Pair{Word: word, Count: count}
		}
	}
}

func reducer(in <-chan Pair) wordcount.WordCount {
	result := make(wordcount.WordCount)
	for pair := range in {
		result[pair.Word] += pair.Count
	}
	return result
}

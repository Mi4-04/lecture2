package main

import (
	"bufio"
	"fmt"
	"lecture2/nlp"
	"lecture2/wordcount"
	"os"
	"path/filepath"
	"sync"
)

type Pair struct {
	Word  string
	Count int
}

const workerCount = 4
const topWordsCount = 100
const readerWorkerCount = 4

func main() {
	files := []string{
		"tom-1.txt",
		"tom-2.txt",
		"tom-3.txt",
		"tom-4.txt",
	}

	fileJobs := make(chan string, len(files))
	lines := make(chan string, readerWorkerCount)
	pairs := make(chan Pair, workerCount)

	go func() {
		for _, file := range files {
			fileJobs <- file
		}
		close(fileJobs)
	}()

	var readerWg sync.WaitGroup
	for i := 0; i < readerWorkerCount; i++ {
		readerWg.Add(1)
		go func() {
			defer readerWg.Done()
			for file := range fileJobs {
				readFileLines(file, lines)
			}
		}()
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
		words := nlp.Tokenize(line)
		counted := make(map[string]int)
		for _, word := range words {
			counted[word]++
		}
		for word, count := range counted {
			out <- Pair{Word: word, Count: count}
		}
	}
}

func reducer(in <-chan Pair) nlp.WordCount {
	result := make(nlp.WordCount)
	for pair := range in {
		result[pair.Word] += pair.Count
	}
	return result
}

func readFileLines(filename string, out chan<- string) {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		fmt.Printf("Ошибка открытия файла %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out <- scanner.Text()
	}
}

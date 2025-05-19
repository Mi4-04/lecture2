package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type WordCount map[string]int

var (
	wordRegex     = regexp.MustCompile(`\p{L}+`)
	prepositions  = map[string]struct{}{"в": {}, "на": {}, "и": {}, "не": {}, "что": {}, "с": {}, "по": {}, "о": {}, "за": {}, "от": {}, "у": {}, "к": {}, "до": {}, "из": {}, "без": {}, "для": {}, "при": {}, "про": {}, "об": {}, "под": {}, "а": {}, "то": {}, "но": {}, "ли": {}, "же": {}, "ни": {}, "ну": {}}
	cyrillicCheck = regexp.MustCompile(`^[а-яА-ЯёЁ]+$`)
)

type Pair struct {
	Word  string
	Count int
}

func main() {
	files := []string{
		"tom-1.txt",
		"tom-2.txt",
		"tom-3.txt",
		"tom-4.txt",
	}

	mapOut := make(chan Pair, 1000)
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			mapper(f, mapOut)
		}(file)
	}

	go func() {
		wg.Wait()
		close(mapOut)
	}()


	words := reduce(mapOut)

	finalCounts := reducer(words)

	printTopWords(finalCounts, 100)
}

func mapper(filename string, out chan<- Pair) {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		fmt.Printf("Ошибка открытия файла %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		

		// фаин фаоут в отдельную горутины с 73 по 86 кол воркеров ограниченно
		words := tokenize(line)
		counted := map[string]int{}

		for _, word := range words {
			if _, ok := prepositions[word]; ok || !isRussian(word) {
				continue
			}
			word = normalize(word)
			counted[word]++
		}

		for word, count := range counted {
			out <- Pair{Word: word, Count: count}
		}
	}
}

// переместить эту логику в reducer и сразу аккумулировать с канала
func reduce(in <-chan Pair) map[string][]int {
	shuffled := make(map[string][]int)

	for pair := range in {
		shuffled[pair.Word] = append(shuffled[pair.Word], pair.Count)
	}

	return shuffled
}

// тут канал арг 
func reducer(shuffled map[string][]int) WordCount {
	result := make(WordCount)
	for word, counts := range shuffled {
		total := 0
		for _, c := range counts {
			total += c
		}
		result[word] = total
	}
	return result
}

func tokenize(line string) []string {
	tokens := wordRegex.FindAllString(line, -1)
	for i, token := range tokens {
		tokens[i] = strings.ToLower(token)
	}
	return tokens
}

func isRussian(word string) bool {
	return cyrillicCheck.MatchString(word)
}

func normalize(word string) string {

	return word
}

func printTopWords(counts WordCount, top int) {
	type kv struct {
		Word  string
		Count int
	}


	//slice is copy func 
	var sorted []kv
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	// срез слайса и все
	fmt.Printf("\nТоп %d слов:\n", top)
	for i := 0; i < top && i < len(sorted); i++ {
		fmt.Printf("%2d. %s — %d\n", i+1, sorted[i].Word, sorted[i].Count)
	}
}

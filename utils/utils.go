package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	WordRegex    = regexp.MustCompile(`\p{L}+`)
	Prepositions = map[string]struct{}{
		"в": {}, "на": {}, "и": {}, "не": {}, "что": {}, "с": {}, "по": {}, "о": {}, "за": {}, "от": {},
		"у": {}, "к": {}, "до": {}, "из": {}, "без": {}, "для": {}, "при": {}, "про": {}, "об": {},
		"под": {}, "а": {}, "то": {}, "но": {}, "ли": {}, "же": {}, "ни": {}, "ну": {},
	}
	CyrillicCheck = regexp.MustCompile(`^[а-яА-ЯёЁ]+$`)
)

func ReadFileLines(filename string, out chan<- string) {
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

func Tokenize(line string) []string {
	tokens := WordRegex.FindAllString(line, -1)
	for i, token := range tokens {
		tokens[i] = strings.ToLower(token)
	}
	return tokens
}

func IsRussian(word string) bool {
	return CyrillicCheck.MatchString(word)
}

func Normalize(word string) string {
	return word
}

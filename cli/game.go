package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)


type game struct {
	words []string
	letters [2]string
	started bool
	nextWordIndex int
	currentWordTemplate string
	currentWord string
	currentWordCorrect string
	nextWordLine int
	wordsNotGuessed []string
	lastAvailableLine int
	wordsGiven int
}

func newGame() *game {
	return &game{
		words: []string{},
		letters: [2]string{"", ""},
		started: false,
		nextWordIndex: 0,
		currentWordTemplate: "",
		currentWord: "",
		currentWordCorrect: "",
		nextWordLine: 0,
		wordsNotGuessed: []string{},
		lastAvailableLine: 0,
		wordsGiven: 0,
	}
}

func (g *game) readWordsFromFile(fp string) {
	fn := fp
	f, err := os.Open(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening wordlist file %s: %s", fn, err.Error())
		os.Exit(1)
	}
	defer f.Close()

	// TODO: Validation - for now, code assumes that the file contains correct data
	i := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		i++
		if i == 1 {
			lineArr := strings.Split(line, ":")
			g.letters = [2]string{lineArr[0], lineArr[1]}
			continue
		}
		g.words = append(g.words, line)
	}
}

func (g *game) randomizeWords() {
	rand.Shuffle(len(g.words), func(i, j int) {
		g.words[i], g.words[j] = g.words[j], g.words[i]
	})
}

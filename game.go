package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

const NOT_STARTED = 0
const CONTINUE_GAME = 1
const GAME_OVER = 2
const INCORRECT_GUESS = 3
const CORRECT_GUESS = 4

type game struct {
	words               []string
	letters             [2]string
	started             bool
	nextWordIndex       int
	currentWordTemplate string
	currentWord         string
	currentWordCorrect  string
	nextWordLine        int
	wordsNotGuessed     []string
	lastAvailableLine   int
	wordsGiven          int
	availableLines      int
}

func newGame() *game {
	return &game{
		words:               []string{},
		letters:             [2]string{"", ""},
		started:             false,
		nextWordIndex:       0,
		currentWordTemplate: "",
		currentWord:         "",
		currentWordCorrect:  "",
		nextWordLine:        0,
		wordsNotGuessed:     []string{},
		lastAvailableLine:   0,
		wordsGiven:          0,
		availableLines:      20,
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

func (g *game) getLeftLetter() string {
	return g.letters[0]
}

func (g *game) getRightLetter() string {
	return g.letters[1]
}

func (g *game) isStarted() bool {
	return g.started
}

func (g *game) stopGame() {
	g.started = false
}

func (g *game) startGame() {
	g.started = true
	g.nextWordIndex = 0
	g.currentWord = ""
	g.nextWordLine = 0
	g.wordsNotGuessed = []string{}
	g.wordsGiven = 0
}

func (g *game) getNumberOfCorrectAnswers() int {
	return g.wordsGiven - len(g.wordsNotGuessed)
}

func (g *game) getNumberOfUsedWords() int {
	return g.wordsGiven
}

func (g *game) getNumberOfAllWords() int {
	return len(g.words)
}

func (g *game) setCurrentWordWithLeftLetter() {
	g.currentWord = strings.Replace(g.currentWordTemplate, "_", g.letters[0], 1)
}

func (g *game) setCurrentWordWithRightLetter() {
	g.currentWord = strings.Replace(g.currentWordTemplate, "_", g.letters[1], 1)
}

func (g *game) getCurrentWord() string {
	return g.currentWord
}

func (g *game) getCurrentLine() int {
	return g.nextWordLine - 1
}

func (g *game) isCurrentWordEmpty() bool {
	return g.currentWord == ""
}

func (g *game) useNewWord() {
	currentWordArr := strings.Split(g.words[g.nextWordIndex], ":")
	g.currentWordTemplate = currentWordArr[0]
	g.currentWord = g.currentWordTemplate
	g.currentWordCorrect = strings.Replace(g.currentWordTemplate, "_", currentWordArr[1], 1)
	g.nextWordIndex++
	g.nextWordLine = 0
}

func (g *game) setAvailableLines(l int) {
	g.availableLines = l
}

func (g *game) iterate() int {
	// If there is no word then take the next one
	if g.isCurrentWordEmpty() {
		g.useNewWord()
	}

	// We need a position that is at the very bottom
	g.lastAvailableLine = g.availableLines - 2 - len(g.wordsNotGuessed)

	if g.lastAvailableLine == 0 || g.nextWordIndex == len(g.words) {
		g.stopGame()
		g.wordsGiven++
		return GAME_OVER
	}

	// If the word is already in the last line
	if g.nextWordLine == g.lastAvailableLine-1 {
		g.wordsGiven++
		if g.currentWord != g.currentWordCorrect {
			g.wordsNotGuessed = append(g.wordsNotGuessed, g.currentWord)
			g.currentWord = ""
			return INCORRECT_GUESS
		} else {
			g.currentWord = ""
			return CORRECT_GUESS
		}
	}

	// Increment the line for the next iteration
	g.nextWordLine++

	// If the word is not at the bottom then just continue
	return CONTINUE_GAME
}

func (g *game) setNextLineToLast() {
	g.nextWordLine = g.lastAvailableLine - 1
}

func (g *game) getLastLine() int {
	return g.lastAvailableLine - 1
}
